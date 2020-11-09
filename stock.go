package main

import (
	"net/http"
	"regexp"
	"sync"
)

var client = &http.Client{Transport: &http.Transport{Proxy: nil}}

type stock interface {
	realtime() realtime
	chart() chart
}

type realtime struct {
	Index   string    `json:"index"`
	Code    string    `json:"code"`
	Name    string    `json:"name"`
	Now     float64   `json:"now"`
	Change  float64   `json:"change"`
	Percent string    `json:"percent"`
	Sell5   []sellbuy `json:"sell5"`
	Buy5    []sellbuy `json:"buy5"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Open    float64   `json:"open"`
	Last    float64   `json:"last"`
	Update  string    `json:"update"`
}

type sellbuy struct {
	Price  float64
	Volume int
}

type chart struct {
	Last float64 `json:"last"`
	Data []point `json:"chart"`
}

type point struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

type suggest struct {
	Index string
	Code  string
	Name  string
	Type  string
}

func initStock(index, code string) (s stock) {
	switch index {
	case "SSE":
		re := regexp.MustCompile(ssePattern)
		if re.MatchString(code) {
			s = &sse{Code: code}
		}
	case "SZSE":
		re := regexp.MustCompile(szsePattern)
		if re.MatchString(code) {
			s = &szse{Code: code}
		}
	}
	return
}

func doGetRealtime(index, code string) realtime {
	s := initStock(index, code)
	return s.realtime()
}

func doGetChart(index, code string) chart {
	s := initStock(index, code)
	return s.chart()
}

func doGetRealtimes(s []stock) []realtime {
	r := make([]realtime, len(s))
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
