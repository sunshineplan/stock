package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sunshineplan/gohttp"
)

const szsePattern = `(00[0-3]|159|300|399)\d{3}`

type szse struct {
	Code     string
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
	if err := gohttp.GetWithClient(
		"http://www.szse.cn/api/market/ssjjhq/getTimeData?marketId=1&code="+s.Code, nil, client).JSON(&result); err != nil {
		log.Println("Failed to get szse:", err)
		return
	}
	if result.Code != "0" {
		log.Println("Data code not equal zero.")
		return
	}
	s.Realtime.Index = "SZSE"
	s.Realtime.Code = s.Code
	s.Realtime.Name = result.Data.Name
	s.Realtime.Now, _ = strconv.ParseFloat(result.Data.Now, 64)
	s.Realtime.Change, _ = strconv.ParseFloat(result.Data.Delta, 64)
	s.Realtime.Percent = result.Data.DeltaPercent + "%"
	s.Realtime.High, _ = strconv.ParseFloat(result.Data.High, 64)
	s.Realtime.Low, _ = strconv.ParseFloat(result.Data.Low, 64)
	s.Realtime.Open, _ = strconv.ParseFloat(result.Data.Open, 64)
	s.Realtime.Last, _ = strconv.ParseFloat(result.Data.Close, 64)
	s.Realtime.Update = result.Data.MarketTime
	var sell5 []sellbuy
	var buy5 []sellbuy
	for i, v := range result.Data.Sellbuy5 {
		price, _ := strconv.ParseFloat(v.Price, 64)
		if i < 5 {
			sell5 = append(sell5, sellbuy{Price: price, Volume: v.Volume})
		} else {
			buy5 = append(buy5, sellbuy{Price: price, Volume: v.Volume})
		}
	}
	s.Realtime.Sell5 = sell5
	s.Realtime.Buy5 = buy5
	s.Chart.Last = s.Realtime.Last
	for _, i := range result.Data.PicUpData {
		y, _ := strconv.ParseFloat(i[1].(string), 64)
		s.Chart.Data = append(s.Chart.Data, point{X: i[0].(string), Y: y})
	}
}

func (s *szse) realtime() realtime {
	s.getRealtime()
	return s.Realtime
}

func (s *szse) chart() chart {
	s.getRealtime()
	return s.Chart
}

func szseSuggest(keyword string) (suggests []suggest) {
	var result []struct{ WordB, Value, Type string }
	if err := gohttp.PostWithClient(
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
