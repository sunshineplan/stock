package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	ssePattern  = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`
	szsePattern = `(00[0-3]|159|300|399)\d{3}`
)

var client = &http.Client{Transport: &http.Transport{Proxy: nil}, Timeout: 2 * time.Second}

type stock interface {
	realtime() map[string]interface{}
	chart() map[string]interface{}
}

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
	url := "http://yunhq.sse.com.cn:32041/v1/sh1/snap/" + s.code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do Request Error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll body Error: %v", err)
		return
	}
	d := simplifiedchinese.GBK.NewDecoder()
	utf8data, err := d.Bytes(data)
	if err != nil {
		log.Printf("Fail to convert gb2312: %v", err)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(utf8data, &jsonData); err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	jd := jsonData.(map[string]interface{})
	snap := jd["snap"].([]interface{})
	s.name = snap[0].(string)
	s.now = snap[5].(float64)
	s.change = snap[6].(float64)
	s.percent = fmt.Sprintf("%.3f", snap[7].(float64)) + "%"
	s.high = snap[3].(float64)
	s.low = snap[4].(float64)
	s.open = snap[2].(float64)
	s.last = snap[1].(float64)
	s.update = fmt.Sprintf("%.0f.%.0f", jd["date"].(float64), jd["time"].(float64))
	var sell5 [][]float64
	for i := 0; i < len(snap[len(snap)-1].([]interface{})); i += 2 {
		sell5 = append(sell5, []float64{snap[len(snap)-1].([]interface{})[i].(float64), snap[len(snap)-1].([]interface{})[i+1].(float64)})
	}
	s.sell5 = sell5
	var buy5 [][]float64
	for i := 0; i < len(snap[len(snap)-2].([]interface{})); i += 2 {
		buy5 = append(buy5, []float64{snap[len(snap)-1].([]interface{})[i].(float64), snap[len(snap)-1].([]interface{})[i+1].(float64)})
	}
	s.buy5 = buy5
}

func (s *sse) getChart() {
	url := "http://yunhq.sse.com.cn:32041/v1/sh1/line/" + s.code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do Request Error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll body Error: %v", err)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("Unmarshal json Error: %v", err)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	req.Header.Set("Referer", "http://www.sse.com.cn/")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do Request Error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll body Error: %v", err)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("Unmarshal json Error: %v", err)
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

type szse struct {
	code      string
	name      string
	now       float64
	change    float64
	percent   string
	sell5     [][]interface{}
	buy5      [][]interface{}
	high      float64
	low       float64
	open      float64
	last      float64
	update    string
	chartData []map[string]interface{}
}

func (s *szse) getRealtime() {
	url := "http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code=" + s.code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do Request Error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll body Error: %v", err)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	if jsonData.(map[string]interface{})["code"] != "0" {
		log.Println("Data code not equal zero.")
		return
	}
	d := jsonData.(map[string]interface{})["data"].(map[string]interface{})
	s.name = d["name"].(string)
	s.now, _ = strconv.ParseFloat(d["now"].(string), 64)
	s.change, _ = strconv.ParseFloat(d["delta"].(string), 64)
	s.percent = d["deltaPercent"].(string) + "%"
	s.high, _ = strconv.ParseFloat(d["high"].(string), 64)
	s.low, _ = strconv.ParseFloat(d["low"].(string), 64)
	s.open, _ = strconv.ParseFloat(d["open"].(string), 64)
	s.last, _ = strconv.ParseFloat(d["close"].(string), 64)
	s.update = d["marketTime"].(string)
	var sell5 [][]interface{}
	var buy5 [][]interface{}
	if d["sellbuy5"] != nil {
		for i, v := range d["sellbuy5"].([]interface{}) {
			if i > 5 {
				sell5 = append(sell5, []interface{}{v.(map[string]interface{})["price"].(string), v.(map[string]interface{})["volume"].(float64)})
			} else {
				buy5 = append(buy5, []interface{}{v.(map[string]interface{})["price"].(string), v.(map[string]interface{})["volume"].(float64)})
			}
		}
	}
	s.sell5 = sell5
	s.buy5 = buy5
	var chart []map[string]interface{}
	for _, v := range d["picupdata"].([]interface{}) {
		chart = append(chart, map[string]interface{}{"x": v.([]interface{})[0].(string), "y": v.([]interface{})[1].(string)})
	}
	s.chartData = chart
}

func (s *szse) realtime() map[string]interface{} {
	s.getRealtime()
	var r = map[string]interface{}{
		"index":   "SZSE",
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

func (s *szse) chart() map[string]interface{} {
	s.getRealtime()
	var r = map[string]interface{}{
		"last":  s.last,
		"chart": s.chartData,
	}
	return r
}

func szseSuggest(keyword string) (r []map[string]interface{}) {
	url := "http://www.szse.cn/api/search/suggest?keyword=" + keyword
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do Request Error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll body Error: %v", err)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	suggest := jsonData.([]interface{})
	re := regexp.MustCompile(szsePattern)
	for _, v := range suggest {
		code := strings.ReplaceAll(strings.ReplaceAll(v.(map[string]interface{})["wordB"].(string), `<span class="keyword">`, ""), "</span>", "")
		if re.MatchString(code) {
			r = append(r, map[string]interface{}{
				"category": "SZSE",
				"code":     code,
				"name":     v.(map[string]interface{})["value"],
				"type":     v.(map[string]interface{})["type"],
			})
		}
	}
	return
}

func initStock(index, code string) (s stock) {
	switch index {
	case "SSE":
		re := regexp.MustCompile(ssePattern)
		if re.MatchString(code) {
			s = &sse{code: code}
		}
	case "SZSE":
		re := regexp.MustCompile(szsePattern)
		if re.MatchString(code) {
			s = &szse{code: code}
		}
	}
	return
}

func doGetRealtime(index, code string) map[string]interface{} {
	s := initStock(index, code)
	return s.realtime()
}

func doGetChart(index, code string) map[string]interface{} {
	s := initStock(index, code)
	return s.chart()
}

func doGetRealtimes(s []stock) []map[string]interface{} {
	r := make([]map[string]interface{}, len(s))
	var wg sync.WaitGroup
	for i, v := range s {
		wg.Add(1)
		go func(i int, s stock) {
			defer wg.Done()
			r[i] = s.realtime()
		}(i, v)
	}
	wg.Wait()
	return r
}
