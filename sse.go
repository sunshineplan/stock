package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/sunshineplan/gohttp"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const ssePattern = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`

type sse struct {
	Code     string
	Name     string
	Realtime realtime
	Chart    chart
}

func (s *sse) getRealtime() {
	resp := gohttp.GetWithClient("http://yunhq.sse.com.cn:32041/v1/sh1/snap/"+s.Code, nil, client)
	if resp.Error != nil {
		log.Println("Failed to get sse realtime:", resp.Error)
		return
	}
	d := simplifiedchinese.GBK.NewDecoder()
	utf8data, err := d.Bytes(resp.Bytes())
	if err != nil {
		log.Println("Fail to convert gb2312:", err)
		return
	}
	var r struct {
		Date int
		Time int
		Snap []interface{}
	}
	if err := json.Unmarshal(utf8data, &r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return
	}
	s.Name = r.Snap[0].(string)
	s.Realtime.now = r.Snap[5].(float64)
	s.Realtime.change = r.Snap[6].(float64)
	s.Realtime.percent = fmt.Sprintf("%.2f", r.Snap[7].(float64)) + "%"
	s.Realtime.high = r.Snap[3].(float64)
	s.Realtime.low = r.Snap[4].(float64)
	s.Realtime.open = r.Snap[2].(float64)
	s.Realtime.last = r.Snap[1].(float64)
	s.Realtime.update = fmt.Sprintf("%d.%d", r.Date, r.Time)
	var sell5, buy5 []sellbuy
	for i := 0; i < 10; i += 2 {
		sell5 = append(sell5,
			sellbuy{r.Snap[len(r.Snap)-1].([]interface{})[i].(float64), int(r.Snap[len(r.Snap)-1].([]interface{})[i+1].(float64))})
		buy5 = append(buy5,
			sellbuy{r.Snap[len(r.Snap)-2].([]interface{})[i].(float64), int(r.Snap[len(r.Snap)-2].([]interface{})[i+1].(float64))})
	}
	if !reflect.DeepEqual(sell5, []sellbuy{{}, {}, {}, {}, {}}) {
		s.Realtime.sell5 = sell5
	}
	if !reflect.DeepEqual(buy5, []sellbuy{{}, {}, {}, {}, {}}) {
		s.Realtime.buy5 = buy5
	}
}

func (s *sse) getChart() {
	var r struct {
		PrevClose float64 `json:"prev_close"`
		Line      [][]interface{}
	}
	if err := gohttp.GetWithClient(
		"http://yunhq.sse.com.cn:32041/v1/sh1/line/"+s.Code, nil, client).JSON(&r); err != nil {
		log.Println("Failed to get sse chart:", err)
		return
	}
	s.Realtime.last = r.PrevClose
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
		s.Chart.data = append(s.Chart.data, point{X: sessions[i], Y: v[0].(float64)})
	}
}

func (s *sse) realtime() map[string]interface{} {
	s.getRealtime()
	return map[string]interface{}{
		"index":   "SSE",
		"code":    s.Code,
		"name":    s.Name,
		"now":     s.Realtime.now,
		"change":  s.Realtime.change,
		"percent": s.Realtime.percent,
		"sell5":   s.Realtime.sell5,
		"buy5":    s.Realtime.buy5,
		"high":    s.Realtime.high,
		"low":     s.Realtime.low,
		"open":    s.Realtime.open,
		"last":    s.Realtime.last,
		"update":  s.Realtime.update,
	}
}

func (s *sse) chart() map[string]interface{} {
	s.getChart()
	return map[string]interface{}{
		"last":  s.Realtime.last,
		"chart": s.Chart.data,
	}
}

func sseSuggest(keyword string) (suggests []suggest) {
	var result struct {
		Data []struct{ Category, Code, Word string }
	}
	url := "http://query.sse.com.cn/search/getPrepareSearchResult.do?search=ycxjs&searchword=" + keyword
	headers := map[string]string{"Referer": "http://www.sse.com.cn/"}
	if err := gohttp.GetWithClient(url, headers, client).JSON(&result); err != nil {
		log.Println("Failed to get sse suggest:", err)
		return
	}
	re := regexp.MustCompile(ssePattern)
	for _, i := range result.Data {
		if re.MatchString(i.Code) {
			suggests = append(suggests, suggest{
				Index: "SSE",
				Code:  i.Code,
				Name:  i.Word,
				Type:  i.Category,
			})
		}
	}
	return
}
