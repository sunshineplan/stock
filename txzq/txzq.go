package txzq

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
)

const ssePattern = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`
const szsePattern = `(00[0-3]|159|300|399)\d{3}`

const api = "http://web.ifzq.gtimg.cn/appstock/app/minute/query?code="
const suggestAPI = "http://smartbox.gtimg.cn/s3/?t=%s&q=%s"

// Timeout specifies a time limit for requests.
var Timeout time.Duration

// SetTimeout sets http client timeout when fetching stocks.
func SetTimeout(duration int) {
	Timeout = time.Duration(duration) * time.Second
}

// TXZQ represents 腾讯证券.
type TXZQ struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (t *TXZQ) get() *TXZQ {
	t.Realtime.Index = t.Index
	t.Realtime.Code = t.Code
	var stk string
	switch strings.ToLower(t.Index) {
	case "sse":
		stk = "sh" + t.Code
	case "szse":
		stk = "sz" + t.Code
	default:
		return t
	}
	resp := gohttp.GetWithClient(api+stk, nil, &http.Client{
		Transport: &http.Transport{Proxy: nil},
		Timeout:   Timeout,
	})
	if resp.Error != nil {
		log.Println("Failed to get txzq:", resp.Error)
		return t
	}
	var r struct{ Data map[string]interface{} }
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return t
	}
	data, ok := r.Data[stk]
	if !ok {
		log.Println("Failed to get this stock:", stk)
		return t
	}
	realtime, ok := data.(map[string]interface{})["qt"].(map[string]interface{})[stk].([]interface{})
	if !ok {
		log.Println("Failed to get this stock realtime")
		return t
	}
	t.Realtime.Name = realtime[1].(string)
	t.Realtime.Now, _ = strconv.ParseFloat(realtime[3].(string), 64)
	t.Realtime.Change, _ = strconv.ParseFloat(realtime[31].(string), 64)
	t.Realtime.Percent = realtime[32].(string) + "%"
	t.Realtime.High, _ = strconv.ParseFloat(realtime[33].(string), 64)
	t.Realtime.Low, _ = strconv.ParseFloat(realtime[34].(string), 64)
	t.Realtime.Open, _ = strconv.ParseFloat(realtime[5].(string), 64)
	t.Realtime.Last, _ = strconv.ParseFloat(realtime[4].(string), 64)
	t.Realtime.Update = realtime[30].(string)
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
		t.Realtime.Buy5 = buy5
		t.Realtime.Sell5 = sell5
	} else {
		t.Realtime.Buy5 = []stock.SellBuy{}
		t.Realtime.Sell5 = []stock.SellBuy{}
	}

	chart, ok := data.(map[string]interface{})["data"].(map[string]interface{})["data"].([]interface{})
	if !ok {
		log.Println("Failed to get this stock chart")
		return t
	}
	t.Chart.Last = t.Realtime.Last
	for _, i := range chart {
		point := strings.Split(i.(string), " ")
		x := point[0][0:2] + ":" + point[0][2:4]
		y, _ := strconv.ParseFloat(point[1], 64)
		t.Chart.Data = append(t.Chart.Data, stock.Point{X: x, Y: y})
	}
	return t
}

// GetRealtime gets the stock's realtime information.
func (t *TXZQ) GetRealtime() stock.Realtime {
	return t.get().Realtime
}

// GetChart gets the stock's chart data.
func (t *TXZQ) GetChart() stock.Chart {
	return t.get().Chart
}

// Suggests returns sse and szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	for _, t := range []string{"gp", "jj"} {
		result := gohttp.GetWithClient(fmt.Sprintf(suggestAPI, t, keyword), nil, &http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   Timeout,
		}).String()
		sse := regexp.MustCompile(ssePattern)
		szse := regexp.MustCompile(szsePattern)
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
		ssePattern,
		func(code string) stock.Stock {
			return &TXZQ{Index: "SSE", Code: code}
		},
		Suggests,
		SetTimeout,
	)
	stock.RegisterStock(
		"szse",
		szsePattern,
		func(code string) stock.Stock {
			return &TXZQ{Index: "SZSE", Code: code}
		},
		func(_ string) []stock.Suggest {
			return nil
		},
		SetTimeout,
	)
}
