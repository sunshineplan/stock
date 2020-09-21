package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sunshineplan/utils/requests"
)

const szsePattern = `(00[0-3]|159|300|399)\d{3}`

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
	var jsonData interface{}
	if err := requests.GetWithClient(
		"http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code="+s.code, nil, client).JSON(&jsonData); err != nil {
		log.Println("Failed to get szse:", err)
		return
	}
	if jsonData.(map[string]interface{})["code"] != "0" {
		log.Println("Data code not equal zero.")
		return
	}
	d := jsonData.(map[string]interface{})["data"].(map[string]interface{})
	s.name = d["name"].(string)
	if d["now"] != nil {
		s.now, _ = strconv.ParseFloat(d["now"].(string), 64)
	}
	s.change, _ = strconv.ParseFloat(d["delta"].(string), 64)
	s.percent = d["deltaPercent"].(string) + "%"
	if d["high"] != nil {
		s.high, _ = strconv.ParseFloat(d["high"].(string), 64)
	}
	if d["low"] != nil {
		s.low, _ = strconv.ParseFloat(d["low"].(string), 64)
	}
	if d["open"] != nil {
		s.open, _ = strconv.ParseFloat(d["open"].(string), 64)
	}
	s.last, _ = strconv.ParseFloat(d["close"].(string), 64)
	s.update = d["marketTime"].(string)
	var sell5 [][]interface{}
	var buy5 [][]interface{}
	if d["sellbuy5"] != nil {
		for i, v := range d["sellbuy5"].([]interface{}) {
			if i < 5 {
				sell5 = append(sell5, []interface{}{v.(map[string]interface{})["price"], v.(map[string]interface{})["volume"]})
			} else {
				buy5 = append(buy5, []interface{}{v.(map[string]interface{})["price"], v.(map[string]interface{})["volume"]})
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
	var jsonData interface{}
	if err := requests.PostWithClient(
		"http://www.szse.cn/api/search/suggest?keyword="+keyword, nil, nil, client).JSON(&jsonData); err != nil {
		log.Println("Failed to get szse suggest:", err)
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
