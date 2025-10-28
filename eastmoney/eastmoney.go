package eastmoney

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sunshineplan/stock"
)

const (
	api          = "https://push2.eastmoney.com/api/qt/stock/get?fltt=2&fields=f11,f12,f13,f14,f15,f16,f17,f18,f19,f20,f31,f32,f33,f34,f35,f36,f37,f38,f39,f40,f43,f44,f45,f46,f58,f60,f169,f170,f531&secid="
	chartAPI     = "https://push2.eastmoney.com/api/qt/stock/trends2/get?iscr=0&fields1=f5,f8&fields2=f53&secid="
	suggestAPI   = "https://searchadapter.eastmoney.com/api/suggest/get?type=14&token=%s&input=%s&count=50"
	suggestToken = "D43BF722C8E33BDC906FB84D85E326E8"
)

var _ stock.Stock = &EastMoney{}

// EastMoney represents 东方财富网.
type EastMoney struct {
	Index    string
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

type eastmoney struct {
	Data struct {
		Now     float64 `json:"F43"`
		High    float64 `json:"F44"`
		Low     float64 `json:"F45"`
		Open    float64 `json:"F46"`
		Name    string  `json:"F58"`
		Last    float64 `json:"F60"`
		Change  float64 `json:"F169"`
		Percent float64 `json:"F170"`

		Buy5Price   float64 `json:"F11"`
		Buy5Volume  int     `json:"F12"`
		Buy4Price   float64 `json:"F13"`
		Buy4Volume  int     `json:"F14"`
		Buy3Price   float64 `json:"F15"`
		Buy3Volume  int     `json:"F16"`
		Buy2Price   float64 `json:"F17"`
		Buy2Volume  int     `json:"F18"`
		Buy1Price   float64 `json:"F19"`
		Buy1Volume  int     `json:"F20"`
		Sell5Price  float64 `json:"F31"`
		Sell5Volume int     `json:"F32"`
		Sell4Price  float64 `json:"F33"`
		Sell4Volume int     `json:"F34"`
		Sell3Price  float64 `json:"F35"`
		Sell3Volume int     `json:"F36"`
		Sell2Price  float64 `json:"F37"`
		Sell2Volume int     `json:"F38"`
		Sell1Price  float64 `json:"F39"`
		Sell1Volume int     `json:"F40"`
	}
}

type eastmoneyChart struct {
	Data struct {
		PreClose float64
		Trends   []string
	}
}

func (s *EastMoney) getRealtime() *EastMoney {
	s.Realtime.Index = s.Index
	s.Realtime.Code = s.Code

	var code string
	switch strings.ToLower(s.Index) {
	case "sse":
		code = "1." + s.Code
	case "szse", "bse":
		code = "0." + s.Code
	default:
		return s
	}

	var r eastmoney
	resp, err := stock.Session.Get(api+code, nil)
	if err != nil {
		log.Println("Failed to get eastmoney realtime:", err)
		return s
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return s
	}
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}

	s.Realtime.Name = r.Data.Name
	s.Realtime.Now = r.Data.Now
	s.Realtime.High = r.Data.High
	s.Realtime.Low = r.Data.Low
	s.Realtime.Open = r.Data.Open
	s.Realtime.Last = r.Data.Last
	s.Realtime.Change = r.Data.Change
	s.Realtime.Percent = fmt.Sprintf("%g%%", r.Data.Percent)
	s.Realtime.Update = time.Now().Format(time.RFC3339)

	s.Realtime.Buy5 = []stock.SellBuy{
		{Price: r.Data.Buy1Price, Volume: r.Data.Buy1Volume},
		{Price: r.Data.Buy2Price, Volume: r.Data.Buy2Volume},
		{Price: r.Data.Buy3Price, Volume: r.Data.Buy3Volume},
		{Price: r.Data.Buy4Price, Volume: r.Data.Buy4Volume},
		{Price: r.Data.Buy5Price, Volume: r.Data.Buy5Volume},
	}
	s.Realtime.Sell5 = []stock.SellBuy{
		{Price: r.Data.Sell1Price, Volume: r.Data.Sell1Volume},
		{Price: r.Data.Sell2Price, Volume: r.Data.Sell2Volume},
		{Price: r.Data.Sell3Price, Volume: r.Data.Sell3Volume},
		{Price: r.Data.Sell4Price, Volume: r.Data.Sell4Volume},
		{Price: r.Data.Sell5Price, Volume: r.Data.Sell5Volume},
	}

	if reflect.DeepEqual(s.Realtime.Sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) &&
		reflect.DeepEqual(s.Realtime.Buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		s.Realtime.Buy5 = []stock.SellBuy{}
		s.Realtime.Sell5 = []stock.SellBuy{}
	}

	return s
}

func (s *EastMoney) getChart() *EastMoney {
	var code string
	switch strings.ToLower(s.Index) {
	case "sse":
		code = "1." + s.Code
	case "szse", "bse":
		code = "0." + s.Code
	default:
		return s
	}

	var r eastmoneyChart
	resp, err := stock.Session.Get(chartAPI+code, nil)
	if err != nil {
		log.Println("Failed to get eastmoney chart:", err)
		return s
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return s
	}
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}

	s.Chart.Last = r.Data.PreClose

	for _, i := range r.Data.Trends {
		data := strings.Split(strings.Split(i, " ")[1], ",")
		x := data[0]
		y, _ := strconv.ParseFloat(data[1], 64)
		s.Chart.Data = append(s.Chart.Data, stock.Point{X: x, Y: y})
	}

	return s
}

// GetRealtime gets the stock's realtime information.
func (s *EastMoney) GetRealtime() stock.Realtime {
	return s.getRealtime().Realtime
}

// GetChart gets the stock's chart data.
func (s *EastMoney) GetChart() stock.Chart {
	return s.getChart().Chart
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
	resp, err := stock.Session.Get(fmt.Sprintf(suggestAPI, suggestToken, keyword), nil)
	if err != nil {
		log.Println("Failed to get eastmoney suggest:", err)
		return
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return
	}
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return
	}

	sse := regexp.MustCompile(stock.SSEPattern)
	szse := regexp.MustCompile(stock.SZSEPattern)
	bse := regexp.MustCompile(stock.BSEPattern)

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
		case "_TB":
			if bse.MatchString(i.Code) {
				suggests = append(suggests,
					stock.Suggest{Index: "BSE", Code: i.Code, Name: i.Name, Type: i.SecurityTypeName})
			}
		}
	}

	return
}
