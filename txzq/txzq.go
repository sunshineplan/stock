package txzq

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sunshineplan/stock"
)

const (
	api        = "http://web.ifzq.gtimg.cn/appstock/app/minute/query?code="
	suggestAPI = "http://smartbox.gtimg.cn/s3/?t=%s&q=%s"
)

var _ stock.Stock = &TXZQ{}

// TXZQ represents 腾讯证券.
type TXZQ struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (txzq *TXZQ) get() *TXZQ {
	txzq.Realtime.Index = txzq.Index
	txzq.Realtime.Code = txzq.Code

	var code string
	switch strings.ToLower(txzq.Index) {
	case "sse":
		code = "sh" + txzq.Code
	case "szse":
		code = "sz" + txzq.Code
	default:
		return txzq
	}

	var r struct{ Data map[string]interface{} }
	if err := stock.Session.Get(api+code, nil).JSON(&r); err != nil {
		log.Println("Failed to get txzq:", err)
		return txzq
	}

	data, ok := r.Data[code]
	if !ok {
		log.Println("Failed to get this stock:", code)
		return txzq
	}

	realtime, ok := data.(map[string]interface{})["qt"].(map[string]interface{})[code].([]interface{})
	if !ok {
		log.Println("Failed to get this stock realtime")
		return txzq
	}

	txzq.Realtime.Name = realtime[1].(string)
	txzq.Realtime.Now, _ = strconv.ParseFloat(realtime[3].(string), 64)
	txzq.Realtime.Change, _ = strconv.ParseFloat(realtime[31].(string), 64)
	txzq.Realtime.Percent = realtime[32].(string) + "%"
	txzq.Realtime.High, _ = strconv.ParseFloat(realtime[33].(string), 64)
	txzq.Realtime.Low, _ = strconv.ParseFloat(realtime[34].(string), 64)
	txzq.Realtime.Open, _ = strconv.ParseFloat(realtime[5].(string), 64)
	txzq.Realtime.Last, _ = strconv.ParseFloat(realtime[4].(string), 64)
	txzq.Realtime.Update = realtime[30].(string)

	buy5 := []stock.SellBuy{}
	sell5 := []stock.SellBuy{}
	for i := 9; i < 19; i += 2 {
		price, _ := strconv.ParseFloat(realtime[i].(string), 64)
		volume, _ := strconv.Atoi(realtime[i+1].(string))
		buy5 = append(buy5, stock.SellBuy{Price: price, Volume: volume})
	}
	for i := 19; i < 29; i += 2 {
		price, _ := strconv.ParseFloat(realtime[i].(string), 64)
		volume, _ := strconv.Atoi(realtime[i+1].(string))
		sell5 = append(sell5, stock.SellBuy{Price: price, Volume: volume})
	}
	if !reflect.DeepEqual(buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) ||
		!reflect.DeepEqual(sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		txzq.Realtime.Buy5 = buy5
		txzq.Realtime.Sell5 = sell5
	} else {
		txzq.Realtime.Buy5 = []stock.SellBuy{}
		txzq.Realtime.Sell5 = []stock.SellBuy{}
	}

	chart, ok := data.(map[string]interface{})["data"].(map[string]interface{})["data"].([]interface{})
	if !ok {
		log.Println("Failed to get this stock chart")
		return txzq
	}
	txzq.Chart.Last = txzq.Realtime.Last
	for _, i := range chart {
		point := strings.Split(i.(string), " ")
		x := point[0][0:2] + ":" + point[0][2:4]
		y, _ := strconv.ParseFloat(point[1], 64)
		if x != "13:00" {
			txzq.Chart.Data = append(txzq.Chart.Data, stock.Point{X: x, Y: y})
		}
	}

	return txzq
}

// GetRealtime gets the stock's realtime information.
func (txzq *TXZQ) GetRealtime() stock.Realtime {
	return txzq.get().Realtime
}

// GetChart gets the stock's chart data.
func (txzq *TXZQ) GetChart() stock.Chart {
	return txzq.get().Chart
}

// Suggests returns sse and szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	for _, t := range []string{"gp", "jj"} {
		result := stock.Session.Get(fmt.Sprintf(suggestAPI, t, keyword), nil).String()
		sse := regexp.MustCompile(stock.SSEPattern)
		szse := regexp.MustCompile(stock.SZSEPattern)
		for _, i := range split(result) {
			name, _ := strconv.Unquote(fmt.Sprintf(`"%s"`, i[2]))
			switch i[0] {
			case "sh":
				if sse.MatchString(i[1]) {
					suggests = append(suggests,
						stock.Suggest{Index: "SSE", Code: i[1], Name: name, Type: i[4]})
				}
			case "sz":
				if szse.MatchString(i[1]) {
					suggests = append(suggests,
						stock.Suggest{Index: "SZSE", Code: i[1], Name: name, Type: i[4]})
				}
			}
		}
	}

	return
}

func split(suggest string) (suggests [][]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in", suggest)
		}
	}()

	s := strings.Split(suggest, `"`)
	if s[1] == "N" {
		return
	}

	for _, i := range strings.Split(s[1], "^") {
		suggests = append(suggests, strings.Split(i, "~"))
	}

	return
}

func init() {
	stock.RegisterStock(
		"sse",
		stock.SSEPattern,
		func(code string) stock.Stock {
			return &TXZQ{Index: "SSE", Code: code}
		},
		Suggests,
	)

	stock.RegisterStock(
		"szse",
		stock.SZSEPattern,
		func(code string) stock.Stock {
			return &TXZQ{Index: "SZSE", Code: code}
		},
		func(_ string) []stock.Suggest {
			return nil
		},
	)
}
