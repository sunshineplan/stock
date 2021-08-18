package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const api = "http://yunhq.sse.com.cn:32041/v1/sh1/snap/"
const chartAPI = "http://yunhq.sse.com.cn:32041/v1/sh1/line/"
const suggestAPI = "http://query.sse.com.cn/search/getPrepareSearchResult.do?search=ycxjs&searchword="

// SSE represents Shanghai Stock Exchange.
type SSE struct {
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (sse *SSE) getRealtime() *SSE {
	sse.Realtime.Index = "SSE"
	sse.Realtime.Code = sse.Code

	resp := stock.Session.Get(api+sse.Code, nil)
	if resp.Error != nil {
		log.Println("Failed to get sse realtime:", resp.Error)
		return sse
	}

	d := simplifiedchinese.GBK.NewDecoder()
	utf8data, err := d.Bytes(resp.Bytes())
	if err != nil {
		log.Println("Fail to convert gb2312:", err)
		return sse
	}

	var r struct {
		Date int
		Time int
		Snap []interface{}
	}
	if err := json.Unmarshal(utf8data, &r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return sse
	}

	sse.Realtime.Name = r.Snap[0].(string)
	sse.Realtime.Now = r.Snap[5].(float64)
	sse.Realtime.Change = r.Snap[6].(float64)
	sse.Realtime.Percent = fmt.Sprintf("%.2f", r.Snap[7].(float64)) + "%"
	sse.Realtime.High = r.Snap[3].(float64)
	sse.Realtime.Low = r.Snap[4].(float64)
	sse.Realtime.Open = r.Snap[2].(float64)
	sse.Realtime.Last = r.Snap[1].(float64)
	sse.Realtime.Update = fmt.Sprintf("%d.%d", r.Date, r.Time)

	sell5 := []stock.SellBuy{}
	buy5 := []stock.SellBuy{}
	for i := 0; i < 10; i += 2 {
		sell5 = append(
			sell5,
			stock.SellBuy{
				Price:  r.Snap[len(r.Snap)-1].([]interface{})[i].(float64),
				Volume: int(r.Snap[len(r.Snap)-1].([]interface{})[i+1].(float64)) / 100,
			})
		buy5 = append(
			buy5,
			stock.SellBuy{
				Price:  r.Snap[len(r.Snap)-2].([]interface{})[i].(float64),
				Volume: int(r.Snap[len(r.Snap)-2].([]interface{})[i+1].(float64)) / 100,
			})
	}
	if !reflect.DeepEqual(sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) ||
		!reflect.DeepEqual(buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		sse.Realtime.Sell5 = sell5
		sse.Realtime.Buy5 = buy5
	} else {
		sse.Realtime.Buy5 = []stock.SellBuy{}
		sse.Realtime.Sell5 = []stock.SellBuy{}
	}

	return sse
}

func (sse *SSE) getChart() *SSE {
	var r struct {
		PrevClose float64 `json:"prev_close"`
		Line      [][]interface{}
	}
	if err := stock.Session.Get(chartAPI+sse.Code, nil).JSON(&r); err != nil {
		log.Println("Failed to get sse chart:", err)
		return sse
	}
	sse.Chart.Last = r.PrevClose

	t := time.Now()
	var sessions []string
	for i := 0; i < 121; i++ {
		sessions = append(
			sessions, time.Date(t.Year(), t.Month(), t.Day(), 9, 30, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	for i := 0; i < 120; i++ {
		sessions = append(
			sessions, time.Date(t.Year(), t.Month(), t.Day(), 13, 1, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	for i, v := range r.Line {
		sse.Chart.Data = append(sse.Chart.Data, stock.Point{X: sessions[i], Y: v[0].(float64)})
	}

	return sse
}

// GetRealtime gets the sse stock's realtime information.
func (sse *SSE) GetRealtime() stock.Realtime {
	return sse.getRealtime().Realtime
}

// GetChart gets the sse stock's chart data.
func (sse *SSE) GetChart() stock.Chart {
	return sse.getChart().Chart
}

// Suggests returns sse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	var r struct {
		Data []struct{ Category, Code, Word string }
	}
	if err := gohttp.Get(suggestAPI+keyword, gohttp.H{"Referer": "http://www.sse.com.cn/"}).JSON(&r); err != nil {
		log.Println("Failed to get sse suggest:", err)
		return
	}

	re := regexp.MustCompile(stock.SSEPattern)
	for _, i := range r.Data {
		if re.MatchString(i.Code) {
			suggests = append(suggests, stock.Suggest{
				Index: "SSE",
				Code:  i.Code,
				Name:  i.Word,
				Type:  i.Category,
			})
		}
	}

	return
}

func init() {
	stock.RegisterStock(
		"sse",
		stock.SSEPattern,
		func(code string) stock.Stock {
			return &SSE{Code: code}
		},
		Suggests,
	)
}
