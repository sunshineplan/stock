package sse

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
)

const (
	api        = "http://yunhq.sse.com.cn:32041/v1/sh1/snap/"
	chartAPI   = "http://yunhq.sse.com.cn:32041/v1/sh1/line/"
	suggestAPI = "http://query.sse.com.cn/search/getPrepareSearchResult.do?search=ycxjs&searchword="
)

var _ stock.Stock = &SSE{}

// SSE represents Shanghai Stock Exchange.
type SSE struct {
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

type sse struct {
	Date int
	Time int
	Snap []any
}

type sseChart struct {
	PrevClose float64 `json:"prev_close"`
	Line      [][]any
}

func (s *SSE) getRealtime() *SSE {
	s.Realtime.Index = "SSE"
	s.Realtime.Code = s.Code

	var r sse
	if err := stock.Session.Get(api+s.Code, nil).JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}

	s.Realtime.Name = r.Snap[0].(string)
	s.Realtime.Now = r.Snap[5].(float64)
	s.Realtime.Change = r.Snap[6].(float64)
	s.Realtime.Percent = fmt.Sprintf("%.2f", r.Snap[7].(float64)) + "%"
	s.Realtime.High = r.Snap[3].(float64)
	s.Realtime.Low = r.Snap[4].(float64)
	s.Realtime.Open = r.Snap[2].(float64)
	s.Realtime.Last = r.Snap[1].(float64)
	s.Realtime.Update = fmt.Sprintf("%d.%d", r.Date, r.Time)

	sell5 := []stock.SellBuy{}
	buy5 := []stock.SellBuy{}
	for i := 0; i < 10; i += 2 {
		sell5 = append(
			sell5,
			stock.SellBuy{
				Price:  r.Snap[len(r.Snap)-1].([]any)[i].(float64),
				Volume: int(r.Snap[len(r.Snap)-1].([]any)[i+1].(float64)) / 100,
			})
		buy5 = append(
			buy5,
			stock.SellBuy{
				Price:  r.Snap[len(r.Snap)-2].([]any)[i].(float64),
				Volume: int(r.Snap[len(r.Snap)-2].([]any)[i+1].(float64)) / 100,
			})
	}
	if !reflect.DeepEqual(sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) ||
		!reflect.DeepEqual(buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		s.Realtime.Sell5 = sell5
		s.Realtime.Buy5 = buy5
	} else {
		s.Realtime.Buy5 = []stock.SellBuy{}
		s.Realtime.Sell5 = []stock.SellBuy{}
	}

	return s
}

func (s *SSE) getChart() *SSE {
	var r sseChart
	if err := stock.Session.Get(chartAPI+s.Code, nil).JSON(&r); err != nil {
		log.Println("Failed to get sse chart:", err)
		return s
	}
	s.Chart.Last = r.PrevClose

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
		s.Chart.Data = append(s.Chart.Data, stock.Point{X: sessions[i], Y: v[0].(float64)})
	}

	return s
}

// GetRealtime gets the sse stock's realtime information.
func (s *SSE) GetRealtime() stock.Realtime {
	return s.getRealtime().Realtime
}

// GetChart gets the sse stock's chart data.
func (s *SSE) GetChart() stock.Chart {
	return s.getChart().Chart
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
