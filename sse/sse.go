package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const ssePattern = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`

// Timeout specifies a time limit for requests.
var Timeout time.Duration

// SetTimeout sets http client timeout when fetching stocks.
func SetTimeout(duration int) {
	Timeout = time.Duration(duration) * time.Second
}

// SSE represents Shanghai Stock Exchange.
type SSE struct {
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

func (s *SSE) getRealtime() *SSE {
	s.Realtime.Index = "SSE"
	s.Realtime.Code = s.Code
	resp := gohttp.GetWithClient(
		"http://yunhq.sse.com.cn:32041/v1/sh1/snap/"+s.Code,
		nil,
		&http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   Timeout,
		})
	if resp.Error != nil {
		log.Println("Failed to get sse realtime:", resp.Error)
		return s
	}
	d := simplifiedchinese.GBK.NewDecoder()
	utf8data, err := d.Bytes(resp.Bytes())
	if err != nil {
		log.Println("Fail to convert gb2312:", err)
		return s
	}
	var r struct {
		Date int
		Time int
		Snap []interface{}
	}
	if err := json.Unmarshal(utf8data, &r); err != nil {
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
	var sell5, buy5 []stock.SellBuy
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
		s.Realtime.Sell5 = sell5
		s.Realtime.Buy5 = buy5
	}
	return s
}

func (s *SSE) getChart() *SSE {
	var r struct {
		PrevClose float64 `json:"prev_close"`
		Line      [][]interface{}
	}
	if err := gohttp.GetWithClient(
		"http://yunhq.sse.com.cn:32041/v1/sh1/line/"+s.Code,
		nil,
		&http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   Timeout,
		}).JSON(&r); err != nil {
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
	var result struct {
		Data []struct{ Category, Code, Word string }
	}
	if err := gohttp.GetWithClient(
		"http://query.sse.com.cn/search/getPrepareSearchResult.do?search=ycxjs&searchword="+keyword,
		gohttp.H{"Referer": "http://www.sse.com.cn/"},
		&http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   Timeout,
		}).JSON(&result); err != nil {
		log.Println("Failed to get sse suggest:", err)
		return
	}
	re := regexp.MustCompile(ssePattern)
	for _, i := range result.Data {
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
		ssePattern,
		func(code string) stock.Stock {
			return &SSE{Code: code}
		},
		Suggests,
		SetTimeout,
	)
}
