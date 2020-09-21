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
	Code     string
	Name     string
	Realtime realtime
	Chart    chart
}

func (s *szse) getRealtime() {
	var result struct {
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
	if err := requests.GetWithClient(
		"http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code="+s.Code, nil, client).JSON(&result); err != nil {
		log.Println("Failed to get szse:", err)
		return
	}
	if result.Code != "0" {
		log.Println("Data code not equal zero.")
		return
	}
	s.Name = result.Data.Name
	s.Realtime.now, _ = strconv.ParseFloat(result.Data.Now, 64)
	s.Realtime.change, _ = strconv.ParseFloat(result.Data.Delta, 64)
	s.Realtime.percent = result.Data.DeltaPercent + "%"
	s.Realtime.high, _ = strconv.ParseFloat(result.Data.High, 64)
	s.Realtime.low, _ = strconv.ParseFloat(result.Data.Low, 64)
	s.Realtime.open, _ = strconv.ParseFloat(result.Data.Open, 64)
	s.Realtime.last, _ = strconv.ParseFloat(result.Data.Close, 64)
	s.Realtime.update = result.Data.MarketTime
	var sell5 [][]interface{}
	var buy5 [][]interface{}
	for i, v := range result.Data.Sellbuy5 {
		if i < 5 {
			sell5 = append(sell5, []interface{}{v.Price, v.Volume})
		} else {
			buy5 = append(buy5, []interface{}{v.Price, v.Volume})
		}
	}
	s.Realtime.sell5 = sell5
	s.Realtime.buy5 = buy5
	for _, i := range result.Data.PicUpData {
		y, _ := strconv.ParseFloat(i[1].(string), 64)
		s.Chart.data = append(s.Chart.data, point{X: i[0].(string), Y: y})
	}
}

func (s *szse) realtime() map[string]interface{} {
	s.getRealtime()
	return map[string]interface{}{
		"index":   "SZSE",
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

func (s *szse) chart() map[string]interface{} {
	s.getRealtime()
	return map[string]interface{}{
		"last":  s.Realtime.last,
		"chart": s.Chart.data,
	}
}

func szseSuggest(keyword string) (suggests []suggest) {
	var result []struct{ WordB, Value, Type string }
	if err := requests.PostWithClient(
		"http://www.szse.cn/api/search/suggest?keyword="+keyword, nil, nil, client).JSON(&result); err != nil {
		log.Println("Failed to get szse suggest:", err)
		return
	}
	re := regexp.MustCompile(szsePattern)
	for _, i := range result {
		if code := strings.ReplaceAll(strings.ReplaceAll(i.WordB, `<span class="keyword">`, ""), "</span>", ""); re.MatchString(code) {
			suggests = append(suggests, suggest{
				Index: "SZSE",
				Code:  code,
				Name:  i.Value,
				Type:  i.Type,
			})
		}
	}
	return
}
