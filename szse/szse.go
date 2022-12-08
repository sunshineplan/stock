package szse

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sunshineplan/stock"
)

const (
	api        = "http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code="
	suggestAPI = "http://www.szse.cn/api/search/suggest?keyword="
)

var _ stock.Stock = &SZSE{}

// SZSE represents Shenzhen Stock Exchange.
type SZSE struct {
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

type szse struct {
	Code string
	Data struct {
		Name         string
		Close        string
		Open         string
		Now          string
		High         string
		Low          string
		Delta        string
		DeltaPercent string
		MarketTime   string
		Sellbuy5     []struct {
			Price  string
			Volume int
		}
		PicUpData [][]any
	}
}

func (s *SZSE) get() *SZSE {
	s.Realtime.Index = "SZSE"
	s.Realtime.Code = s.Code

	var r szse
	if err := stock.Session.Get(api+s.Code, nil).JSON(&r); err != nil {
		log.Println("Failed to get szse:", err)
		return s
	}
	if r.Code != "0" {
		log.Println("Data code not equal zero.")
		return s
	}

	s.Realtime.Name = r.Data.Name
	s.Realtime.Now, _ = strconv.ParseFloat(r.Data.Now, 64)
	s.Realtime.Change, _ = strconv.ParseFloat(r.Data.Delta, 64)
	s.Realtime.Percent = r.Data.DeltaPercent + "%"
	s.Realtime.High, _ = strconv.ParseFloat(r.Data.High, 64)
	s.Realtime.Low, _ = strconv.ParseFloat(r.Data.Low, 64)
	s.Realtime.Open, _ = strconv.ParseFloat(r.Data.Open, 64)
	s.Realtime.Last, _ = strconv.ParseFloat(r.Data.Close, 64)
	s.Realtime.Update = r.Data.MarketTime

	sell5 := []stock.SellBuy{}
	buy5 := []stock.SellBuy{}
	for i, v := range r.Data.Sellbuy5 {
		price, _ := strconv.ParseFloat(v.Price, 64)
		if i < 5 {
			sell5 = append(sell5, stock.SellBuy{Price: price, Volume: v.Volume})
		} else {
			buy5 = append(buy5, stock.SellBuy{Price: price, Volume: v.Volume})
		}
	}
	s.Realtime.Sell5 = sell5
	s.Realtime.Buy5 = buy5

	s.Chart.Last = s.Realtime.Last

	for _, i := range r.Data.PicUpData {
		y, _ := strconv.ParseFloat(i[1].(string), 64)
		s.Chart.Data = append(s.Chart.Data, stock.Point{X: i[0].(string), Y: y})
	}

	return s
}

// GetRealtime gets the szse stock's realtime information.
func (s *SZSE) GetRealtime() stock.Realtime {
	return s.get().Realtime
}

// GetChart gets the szse stock's chart data.
func (s *SZSE) GetChart() stock.Chart {
	return s.get().Chart
}

// Suggests returns szse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	var r []struct{ WordB, Value, Type string }
	if err := stock.Session.Post(suggestAPI+keyword, nil, nil).JSON(&r); err != nil {
		log.Println("Failed to get szse suggest:", err)
		return
	}

	re := regexp.MustCompile(stock.SZSEPattern)
	for _, i := range r {
		if code := strings.ReplaceAll(strings.ReplaceAll(i.WordB, `<span class="keyword">`, ""), "</span>", ""); re.MatchString(code) {
			suggests = append(suggests, stock.Suggest{
				Index: "SZSE",
				Code:  code,
				Name:  i.Value,
				Type:  i.Type,
			})
		}
	}

	return
}

func init() {
	stock.RegisterStock(
		"szse",
		stock.SZSEPattern,
		func(code string) stock.Stock {
			return &SZSE{Code: code}
		},
		Suggests,
	)
}
