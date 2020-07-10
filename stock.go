package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type stock interface {
	realtime() map[string]interface{}
	chart() map[string]interface{}
}

type chartPoint struct {
	x string
	y float64
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
	chartData []chartPoint
}

func (s *sse) getRealtime() {
	url := "http://yunhq.sse.com.cn:32041/v1/sh1/snap/" + s.code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
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
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	d := jsonData.(map[string]interface{})
	snap := d["snap"].([]interface{})
	s.name = strings.TrimSpace(snap[0].(string))
	s.now = snap[5].(float64)
	s.change = snap[6].(float64)
	s.percent = snap[7].(string) + "%"
	s.high = snap[3].(float64)
	s.low = snap[4].(float64)
	s.open = snap[2].(float64)
	s.last = snap[1].(float64)
	s.update = d["date"].(string) + "." + d["time"].(string)
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
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
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
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	line := jsonData.(map[string]interface{})["line"].([]interface{})
	t := time.Now()
	var sessions []string
	for i := 0; i < 121; i++ {
		sessions = append(sessions, time.Date(t.Year(), t.Month(), t.Day(), 9, 30, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	for i := 0; i < 120; i++ {
		sessions = append(sessions, time.Date(t.Year(), t.Month(), t.Day(), 13, 1, 0, 0, time.Local).Add(time.Duration(i)*time.Minute).Format("15:04"))
	}
	var chart []chartPoint
	for i, v := range line {
		chart = append(chart, chartPoint{x: sessions[i], y: v.([]interface{})[0].(float64)})
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

type szse struct {
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
	chartData []chartPoint
}

func (s *szse) getRealtime() {
	url := "http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code=" + s.code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("New Request Error: %v", err)
		return
	}
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
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
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Printf("Unmarshal json Error: %v", err)
		return
	}
	if jsonData.(map[string]interface{})["code"] != 0 {
		log.Printf("Data code not equal zero: %v", err)
		return
	}
	d := jsonData.(map[string]interface{})["data"].(map[string]interface{})
	s.name = d["name"].(string)
	s.now = d["now"].(float64)
	s.change = d["delta"].(float64)
	s.percent = d["deltaPercent"].(string) + "%"
	s.high = d["high"].(float64)
	s.low = d["low"].(float64)
	s.open = d["open"].(float64)
	s.last = d["close"].(float64)
	s.update = d["marketTime"].(string)
	var sell5 [][]float64
	var buy5 [][]float64
	for i, v := range d["sellbuy5"].([]interface{}) {
		if i > 5 {
			sell5 = append(sell5, []float64{v.(map[string]interface{})["price"].(float64), v.(map[string]interface{})["volume"].(float64)})
		} else {
			buy5 = append(buy5, []float64{v.(map[string]interface{})["price"].(float64), v.(map[string]interface{})["volume"].(float64)})
		}
	}
	s.sell5 = sell5
	s.buy5 = buy5
	var chart []chartPoint
	for _, v := range d["picupdata"].([]interface{}) {
		chart = append(chart, chartPoint{x: v.([]interface{})[0].(string), y: v.([]interface{})[1].(float64)})
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

func initStock(index, code string) (s stock) {
	switch index {
	case "SSE":
		re := regexp.MustCompile(`000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`)
		if re.MatchString(code) {
			s = &sse{code: code}
		}
	case "SZSE":
		re := regexp.MustCompile(`(00[0-3]|159|300|399)\d{3}`)
		if re.MatchString(code) {
			s = &szse{code: code}
		}
	}
	return
}

func doGetRealtimes(s []stock) []map[string]interface{} {
	var wg sync.WaitGroup
	for _, i := range s {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i.realtime()
		}()
	}
	wg.Wait()
	return nil
}

func doGetRealtime(index, code string) map[string]interface{} {
	s := initStock(index, code)
	return s.realtime()
}

func doGetChart(index, code string) map[string]interface{} {
	s := initStock(index, code)
	return s.chart()
}
