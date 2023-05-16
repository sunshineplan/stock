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
	api        = "https://web.ifzq.gtimg.cn/appstock/app/minute/query?code="
	suggestAPI = "https://smartbox.gtimg.cn/s3/?t=%s&q=%s"
)

var _ stock.Stock = &TXZQ{}

// TXZQ represents 腾讯证券.
type TXZQ struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (s *TXZQ) get() *TXZQ {
	s.Realtime.Index = s.Index
	s.Realtime.Code = s.Code

	var code string
	switch strings.ToLower(s.Index) {
	case "sse":
		code = "sh" + s.Code
	case "szse":
		code = "sz" + s.Code
	default:
		return s
	}

	var r struct{ Data map[string]any }
	resp, err := stock.Session.Get(api+code, nil)
	if err != nil {
		log.Println("Failed to get txzq:", err)
		return s
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return s
	}
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}

	data, ok := r.Data[code]
	if !ok {
		log.Println("Failed to get this stock:", code)
		return s
	}

	realtime, ok := data.(map[string]any)["qt"].(map[string]any)[code].([]any)
	if !ok {
		log.Println("Failed to get this stock realtime")
		return s
	}

	s.Realtime.Name = realtime[1].(string)
	s.Realtime.Now, _ = strconv.ParseFloat(realtime[3].(string), 64)
	s.Realtime.Change, _ = strconv.ParseFloat(realtime[31].(string), 64)
	s.Realtime.Percent = realtime[32].(string) + "%"
	s.Realtime.High, _ = strconv.ParseFloat(realtime[33].(string), 64)
	s.Realtime.Low, _ = strconv.ParseFloat(realtime[34].(string), 64)
	s.Realtime.Open, _ = strconv.ParseFloat(realtime[5].(string), 64)
	s.Realtime.Last, _ = strconv.ParseFloat(realtime[4].(string), 64)
	s.Realtime.Update = realtime[30].(string)

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
		s.Realtime.Buy5 = buy5
		s.Realtime.Sell5 = sell5
	} else {
		s.Realtime.Buy5 = []stock.SellBuy{}
		s.Realtime.Sell5 = []stock.SellBuy{}
	}

	chart, ok := data.(map[string]any)["data"].(map[string]any)["data"].([]any)
	if !ok {
		log.Println("Failed to get this stock chart")
		return s
	}
	s.Chart.Last = s.Realtime.Last
	for _, i := range chart {
		point := strings.Split(i.(string), " ")
		x := point[0][0:2] + ":" + point[0][2:4]
		y, _ := strconv.ParseFloat(point[1], 64)
		if x != "13:00" {
			s.Chart.Data = append(s.Chart.Data, stock.Point{X: x, Y: y})
		}
	}

	return s
}

// GetRealtime gets the stock's realtime information.
func (s *TXZQ) GetRealtime() stock.Realtime {
	return s.get().Realtime
}

// GetChart gets the stock's chart data.
func (s *TXZQ) GetChart() stock.Chart {
	return s.get().Chart
}

// Suggests returns sse and szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	for _, t := range []string{"gp", "jj"} {
		resp, err := stock.Session.Get(fmt.Sprintf(suggestAPI, t, keyword), nil)
		if err != nil {
			log.Println("Failed to get txzq suggest:", err)
			continue
		} else if resp.StatusCode != 200 {
			log.Println("Bad status code:", resp.StatusCode)
			continue
		}
		sse := regexp.MustCompile(stock.SSEPattern)
		szse := regexp.MustCompile(stock.SZSEPattern)
		for _, i := range split(resp.String()) {
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
