package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/sunshineplan/utils/requests"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const ssePattern = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`

type sse struct {
	code      string
	name      string
	now       float64
	change    float64
	percent   string
	sell5     [][]float64
	buy5      [][]float64
	high      float64
	low       float64
	open      float64
	last      float64
	update    string
	chartData []map[string]interface{}
}

func (s *sse) getRealtime() {
	resp := requests.GetWithClient("http://yunhq.sse.com.cn:32041/v1/sh1/snap/"+s.code, nil, client)
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
	var jsonData interface{}
	if err := json.Unmarshal(utf8data, &jsonData); err != nil {
		log.Println("Unmarshal json Error:", err)
		return
	}
	jd := jsonData.(map[string]interface{})
	snap := jd["snap"].([]interface{})
	s.name = snap[0].(string)
	s.now = snap[5].(float64)
	s.change = snap[6].(float64)
	s.percent = fmt.Sprintf("%.2f", snap[7].(float64)) + "%"
	s.high = snap[3].(float64)
	s.low = snap[4].(float64)
	s.open = snap[2].(float64)
	s.last = snap[1].(float64)
	s.update = fmt.Sprintf("%.0f.%.0f", jd["date"].(float64), jd["time"].(float64))
	var sell5, buy5 [][]float64
	for i := 0; i < len(snap[len(snap)-1].([]interface{})); i += 2 {
		sell5 = append(sell5, []float64{snap[len(snap)-1].([]interface{})[i].(float64), snap[len(snap)-1].([]interface{})[i+1].(float64)})
		buy5 = append(buy5, []float64{snap[len(snap)-2].([]interface{})[i].(float64), snap[len(snap)-2].([]interface{})[i+1].(float64)})
	}
	s.sell5 = sell5
	s.buy5 = buy5
}

func (s *sse) getChart() {
	var jsonData interface{}
	if err := requests.GetWithClient(
		"http://yunhq.sse.com.cn:32041/v1/sh1/line/"+s.code, nil, client).JSON(&jsonData); err != nil {
		log.Println("Failed to get sse chart:", err)
		return
	}
	s.last = jsonData.(map[string]interface{})["prev_close"].(float64)
	line := jsonData.(map[string]interface{})["line"].([]interface{})
	t := time.Now()
	var sessions []string
	for i := 0; i < 121; i++ {
		sessions = append(sessions, time.Date(t.Year(), t.Month(), t.Day(), 9, 30, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	for i := 0; i < 120; i++ {
		sessions = append(sessions, time.Date(t.Year(), t.Month(), t.Day(), 13, 1, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	var chart []map[string]interface{}
	for i, v := range line {
		chart = append(chart, map[string]interface{}{"x": sessions[i], "y": v.([]interface{})[0].(float64)})
	}
	s.chartData = chart
}

func (s *sse) realtime() map[string]interface{} {
	s.getRealtime()
	var r = map[string]interface{}{
		"index":   "SSE",
		"code":    s.code,
		"name":    s.name,
		"now":     s.now,
		"change":  s.change,
		"percent": s.percent,
		"sell5":   s.sell5,
		"buy5":    s.buy5,
		"high":    s.high,
		"low":     s.low,
		"open":    s.open,
		"last":    s.last,
		"update":  s.update,
	}
	return r
}

func (s *sse) chart() map[string]interface{} {
	s.getChart()
	var r = map[string]interface{}{
		"last":  s.last,
		"chart": s.chartData,
	}
	return r
}

func sseSuggest(keyword string) (r []map[string]interface{}) {
	url := "http://query.sse.com.cn/search/getPrepareSearchResult.do?search=ycxjs&searchword=" + keyword
	headers := map[string]string{"Referer": "http://www.sse.com.cn/"}
	var jsonData interface{}
	if err := requests.GetWithClient(url, headers, client).JSON(&jsonData); err != nil {
		log.Println("Failed to get sse suggest:", err)
		return
	}
	suggest := jsonData.(map[string]interface{})["data"].([]interface{})
	re := regexp.MustCompile(ssePattern)
	for _, v := range suggest {
		if re.MatchString(v.(map[string]interface{})["CODE"].(string)) {
			r = append(r, map[string]interface{}{
				"category": "SSE",
				"code":     v.(map[string]interface{})["CODE"],
				"name":     v.(map[string]interface{})["WORD"],
				"type":     v.(map[string]interface{})["CATEGORY"],
			})
		}
	}
	return
}
