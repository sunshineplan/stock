package eastmoney

import (
	"fmt"
	"log"
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

var s = gohttp.NewSession()

// SetTimeout sets http client timeout when fetching stocks.
func SetTimeout(duration int) {
	s.SetTimeout(time.Duration(duration) * time.Second)
}

// EastMoney represents 东方财富网.
type EastMoney struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (eastmoney *EastMoney) getRealtime() *EastMoney {
	eastmoney.Realtime.Index = eastmoney.Index
	eastmoney.Realtime.Code = eastmoney.Code

	var code string
	switch strings.ToLower(eastmoney.Index) {
	case "sse":
		code = "1." + eastmoney.Code
	case "szse":
		code = "0." + eastmoney.Code
	default:
		return eastmoney
	}

	var r struct {
		Data struct {
			F11  float64 // buy5 price
			F12  int     // buy5 volume
			F13  float64 // buy4 price
			F14  int     // buy4 volume
			F15  float64 // buy3 price
			F16  int     // buy3 volume
			F17  float64 // buy2 price
			F18  int     // buy2 volume
			F19  float64 // buy1 price
			F20  int     // buy1 volume
			F31  float64 // sell5 price
			F32  int     // sell5 volume
			F33  float64 // sell4 price
			F34  int     // sell4 volume
			F35  float64 // sell3 price
			F36  int     // sell3 volume
			F37  float64 // sell2 price
			F38  int     // sell2 volume
			F39  float64 // sell1 price
			F40  int     // sell1 volume
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
	if err := s.Get(api+code, nil).JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return eastmoney
	}

	eastmoney.Realtime.Name = r.Data.F58
	eastmoney.Realtime.Now = r.Data.F43
	eastmoney.Realtime.High = r.Data.F44
	eastmoney.Realtime.Low = r.Data.F45
	eastmoney.Realtime.Open = r.Data.F46
	eastmoney.Realtime.Last = r.Data.F60
	eastmoney.Realtime.Change = r.Data.F169
	eastmoney.Realtime.Percent = fmt.Sprintf("%g%%", r.Data.F170)
	eastmoney.Realtime.Update = time.Now().Format(time.RFC3339)

	eastmoney.Realtime.Buy5 = []stock.SellBuy{
		{Price: r.Data.F19, Volume: r.Data.F20},
		{Price: r.Data.F17, Volume: r.Data.F18},
		{Price: r.Data.F15, Volume: r.Data.F16},
		{Price: r.Data.F13, Volume: r.Data.F14},
		{Price: r.Data.F11, Volume: r.Data.F12},
	}
	eastmoney.Realtime.Sell5 = []stock.SellBuy{
		{Price: r.Data.F39, Volume: r.Data.F40},
		{Price: r.Data.F37, Volume: r.Data.F38},
		{Price: r.Data.F35, Volume: r.Data.F36},
		{Price: r.Data.F33, Volume: r.Data.F34},
		{Price: r.Data.F31, Volume: r.Data.F32},
	}

	if reflect.DeepEqual(eastmoney.Realtime.Sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) &&
		reflect.DeepEqual(eastmoney.Realtime.Buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		eastmoney.Realtime.Buy5 = []stock.SellBuy{}
		eastmoney.Realtime.Sell5 = []stock.SellBuy{}
	}

	return eastmoney
}

func (eastmoney *EastMoney) getChart() *EastMoney {
	var code string
	switch strings.ToLower(eastmoney.Index) {
	case "sse":
		code = "1." + eastmoney.Code
	case "szse":
		code = "0." + eastmoney.Code
	default:
		return eastmoney
	}

	var r struct {
		Data struct {
			PreClose float64
			Trends   []string
		}
	}
	if err := s.Get(chartAPI+code, nil).JSON(&r); err != nil {
		log.Println("Failed to get eastmoney chart:", err)
		return eastmoney
	}

	eastmoney.Chart.Last = r.Data.PreClose

	for _, i := range r.Data.Trends {
		data := strings.Split(strings.Split(i, " ")[1], ",")
		x := data[0]
		y, _ := strconv.ParseFloat(data[1], 64)
		eastmoney.Chart.Data = append(eastmoney.Chart.Data, stock.Point{X: x, Y: y})
	}

	return eastmoney
}

// GetRealtime gets the stock's realtime information.
func (eastmoney *EastMoney) GetRealtime() stock.Realtime {
	return eastmoney.getRealtime().Realtime
}

// GetChart gets the stock's chart data.
func (eastmoney *EastMoney) GetChart() stock.Chart {
	return eastmoney.getChart().Chart
}

// Suggests returns sse and szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	var r struct {
		QuotationCodeTable struct {
			Data []struct {
				Code             string
				Name             string
				MarketType       string
				SecurityTypeName string
			}
		}
	}
	if err := gohttp.Get(fmt.Sprintf(suggestAPI, suggestToken, keyword), nil).JSON(&r); err != nil {
		log.Println("Failed to get eastmoney suggest:", err)
		return
	}

	sse := regexp.MustCompile(ssePattern)
	szse := regexp.MustCompile(szsePattern)

	for _, i := range r.QuotationCodeTable.Data {
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
