package eastmoney

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

const api = "http://push2.eastmoney.com/api/qt/stock/get?fltt=2&fields=f11,f12,f13,f14,f15,f16,f17,f18,f19,f20,f31,f32,f33,f34,f35,f36,f37,f38,f39,f40,f43,f44,f45,f46,f58,f60,f169,f170,f531&secid="
const chartAPI = "http://push2.eastmoney.com/api/qt/stock/trends2/get?iscr=0&fields1=f5,f8&fields2=f53&secid="
const suggestAPI = "http://searchapi.eastmoney.com/api/suggest/get?type=14&token=%s&input=%s"
const suggestToken = "D43BF722C8E33BDC906FB84D85E326E8"

// Timeout specifies a time limit for requeste.
var Timeout time.Duration

// SetTimeout sets http client timeout when fetching stocke.
func SetTimeout(duration int) {
	Timeout = time.Duration(duration) * time.Second
}

// EastMoney represents 东方财富网.
type EastMoney struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (e *EastMoney) getRealtime() *EastMoney {
	e.Realtime.Index = e.Index
	e.Realtime.Code = e.Code

	var stk string
	switch strings.ToLower(e.Index) {
	case "sse":
		stk = "1." + e.Code
	case "szse":
		stk = "0." + e.Code
	default:
		return e
	}

	var r struct {
		Data struct {
			F11  float64 // buy1 price
			F12  int     // buy1 volume
			F13  float64 // buy2 price
			F14  int     // buy2 volume
			F15  float64 // buy3 price
			F16  int     // buy3 volume
			F17  float64 // buy4 price
			F18  int     // buy4 volume
			F19  float64 // buy5 price
			F20  int     // buy5 volume
			F31  float64 // sell1 price
			F32  int     // sell1 volume
			F33  float64 // sell2 price
			F34  int     // sell2 volume
			F35  float64 // sell3 price
			F36  int     // sell3 volume
			F37  float64 // sell4 price
			F38  int     // sell4 volume
			F39  float64 // sell5 price
			F40  int     // sell5 volume
			F43  float64 // now
			F44  float64 // high
			F45  float64 // low
			F46  float64 // open
			F58  string  // name
			F60  float64 // last
			F169 float64 // change
			F170 float64 // percent
		}
	}
	if err := gohttp.GetWithClient(api+stk, nil, &http.Client{
		Transport: &http.Transport{Proxy: nil},
		Timeout:   Timeout,
	}).JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return e
	}

	e.Realtime.Name = r.Data.F58
	e.Realtime.Now = r.Data.F43
	e.Realtime.High = r.Data.F44
	e.Realtime.Low = r.Data.F45
	e.Realtime.Open = r.Data.F46
	e.Realtime.Last = r.Data.F60
	e.Realtime.Change = r.Data.F169
	e.Realtime.Percent = fmt.Sprintf("%g%%", r.Data.F170)
	e.Realtime.Update = time.Now().Format(time.RFC3339)

	e.Realtime.Buy5 = []stock.SellBuy{
		{Price: r.Data.F11, Volume: r.Data.F12},
		{Price: r.Data.F13, Volume: r.Data.F14},
		{Price: r.Data.F15, Volume: r.Data.F16},
		{Price: r.Data.F17, Volume: r.Data.F18},
		{Price: r.Data.F19, Volume: r.Data.F20},
	}
	e.Realtime.Sell5 = []stock.SellBuy{
		{Price: r.Data.F31, Volume: r.Data.F32},
		{Price: r.Data.F33, Volume: r.Data.F34},
		{Price: r.Data.F35, Volume: r.Data.F36},
		{Price: r.Data.F37, Volume: r.Data.F38},
		{Price: r.Data.F39, Volume: r.Data.F40},
	}

	if reflect.DeepEqual(e.Realtime.Sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) &&
		reflect.DeepEqual(e.Realtime.Buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		e.Realtime.Buy5 = []stock.SellBuy{}
		e.Realtime.Sell5 = []stock.SellBuy{}
	}

	return e
}

func (e *EastMoney) getChart() *EastMoney {
	var stk string
	switch strings.ToLower(e.Index) {
	case "sse":
		stk = "1." + e.Code
	case "szse":
		stk = "0." + e.Code
	default:
		return e
	}

	var r struct {
		Data struct {
			PreClose float64
			Trends   []string
		}
	}
	if err := gohttp.GetWithClient(chartAPI+stk, nil, &http.Client{
		Transport: &http.Transport{Proxy: nil},
		Timeout:   Timeout,
	}).JSON(&r); err != nil {
		log.Println("Failed to get eastmoney chart:", err)
		return e
	}

	e.Chart.Last = r.Data.PreClose

	for _, i := range r.Data.Trends {
		data := strings.Split(strings.Split(i, " ")[1], ",")
		x := data[0]
		y, _ := strconv.ParseFloat(data[1], 64)
		e.Chart.Data = append(e.Chart.Data, stock.Point{X: x, Y: y})
	}

	return e
}

// GetRealtime gets the stock's realtime information.
func (e *EastMoney) GetRealtime() stock.Realtime {
	return e.getRealtime().Realtime
}

// GetChart gets the stock's chart data.
func (e *EastMoney) GetChart() stock.Chart {
	return e.getChart().Chart
}

// Suggests returns sse and szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	var result struct {
		QuotationCodeTable struct {
			Data []struct {
				Code             string
				Name             string
				MarketType       string
				SecurityTypeName string
			}
		}
	}
	if err := gohttp.GetWithClient(fmt.Sprintf(suggestAPI, suggestToken, keyword), nil, &http.Client{
		Transport: &http.Transport{Proxy: nil},
		Timeout:   Timeout,
	}).JSON(&result); err != nil {
		log.Println("Failed to get eastmoney suggest:", err)
		return
	}

	sse := regexp.MustCompile(ssePattern)
	szse := regexp.MustCompile(szsePattern)

	for _, i := range result.QuotationCodeTable.Data {
		switch i.MarketType {
		case "1":
			if sse.MatchString(i.Code) {
				suggests = append(suggests,
					stock.Suggest{Index: "SSE", Code: i.Code, Name: i.Name, Type: i.SecurityTypeName})
			}
		case "2":
			if szse.MatchString(i.Code) {
				suggests = append(suggests,
					stock.Suggest{Index: "SZSE", Code: i.Code, Name: i.Name, Type: i.SecurityTypeName})
			}
		}
	}

	return
}

func init() {
	stock.RegisterStock(
		"sse",
		ssePattern,
		func(code string) stock.Stock {
			return &EastMoney{Index: "SSE", Code: code}
		},
		Suggests,
		SetTimeout,
	)
	stock.RegisterStock(
		"szse",
		szsePattern,
		func(code string) stock.Stock {
			return &EastMoney{Index: "SZSE", Code: code}
		},
		func(_ string) []stock.Suggest {
			return nil
		},
		SetTimeout,
	)
}
