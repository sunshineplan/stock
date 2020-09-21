package main

import (
	"net/http"
	"regexp"
	"sync"
	"time"
)

var client = &http.Client{Transport: &http.Transport{Proxy: nil}, Timeout: 2 * time.Second}

type stock interface {
	realtime() map[string]interface{}
	chart() map[string]interface{}
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
