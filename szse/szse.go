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

func (szse *SZSE) get() *SZSE {
	szse.Realtime.Index = "SZSE"
	szse.Realtime.Code = szse.Code

	var r struct {
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
			PicUpData [][]interface{}
		}
	}
	if err := stock.Session.Get(api+szse.Code, nil).JSON(&r); err != nil {
		log.Println("Failed to get szse:", err)
		return szse
	}
	if r.Code != "0" {
		log.Println("Data code not equal zero.")
		return szse
	}

	szse.Realtime.Name = r.Data.Name
	szse.Realtime.Now, _ = strconv.ParseFloat(r.Data.Now, 64)
	szse.Realtime.Change, _ = strconv.ParseFloat(r.Data.Delta, 64)
	szse.Realtime.Percent = r.Data.DeltaPercent + "%"
	szse.Realtime.High, _ = strconv.ParseFloat(r.Data.High, 64)
	szse.Realtime.Low, _ = strconv.ParseFloat(r.Data.Low, 64)
	szse.Realtime.Open, _ = strconv.ParseFloat(r.Data.Open, 64)
	szse.Realtime.Last, _ = strconv.ParseFloat(r.Data.Close, 64)
	szse.Realtime.Update = r.Data.MarketTime

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
	szse.Realtime.Sell5 = sell5
	szse.Realtime.Buy5 = buy5

	szse.Chart.Last = szse.Realtime.Last

	for _, i := range r.Data.PicUpData {
		y, _ := strconv.ParseFloat(i[1].(string), 64)
		szse.Chart.Data = append(szse.Chart.Data, stock.Point{X: i[0].(string), Y: y})
	}

	return szse
}

// GetRealtime gets the szse stock's realtime information.
func (szse *SZSE) GetRealtime() stock.Realtime {
	return szse.get().Realtime
}

// GetChart gets the szse stock's chart data.
func (szse *SZSE) GetChart() stock.Chart {
	return szse.get().Chart
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
