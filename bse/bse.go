package bse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/stock"
)

const (
	api        = "https://www.bse.cn/nqhqController/nqhq.do?xxfcbj=2&zqdm="
	chartAPI   = "https://www.bse.cn/companyEchartsController/getTimeSharingChart/list/%s.do"
	suggestAPI = "https://www.bse.cn/nqxxController/getBSECode.do"
)

var _ stock.Stock = &BSE{}

// BSE represents Beijing Stock Exchange.
type BSE struct {
	Code     string
	Realtime stock.Realtime
	Chart    stock.Chart
}

type bse struct {
	Content []struct {
		Name    string  `json:"hqzqjc"`
		Now     float64 `json:"hqzjcj"`
		High    float64 `json:"hqzgcj"`
		Low     float64 `json:"hqzdcj"`
		Open    float64 `json:"hqjrkp"`
		Last    float64 `json:"hqzrsp"`
		Change  float64 `json:"hqzd"`
		Percent float64 `json:"hqzdf"`

		Buy1Price   float64 `json:"hqbjw1"`
		Buy1Volume  int     `json:"hqbsl1"`
		Buy2Price   float64 `json:"hqbjw2"`
		Buy2Volume  int     `json:"hqbsl2"`
		Buy3Price   float64 `json:"hqbjw3"`
		Buy3Volume  int     `json:"hqbsl3"`
		Buy4Price   float64 `json:"hqbjw4"`
		Buy4Volume  int     `json:"hqbsl4"`
		Buy5Price   float64 `json:"hqbjw5"`
		Buy5Volume  int     `json:"hqbsl5"`
		Sell1Price  float64 `json:"hqsjw1"`
		Sell1Volume int     `json:"hqssl1"`
		Sell2Price  float64 `json:"hqsjw2"`
		Sell2Volume int     `json:"hqssl2"`
		Sell3Price  float64 `json:"hqsjw3"`
		Sell3Volume int     `json:"hqssl3"`
		Sell4Price  float64 `json:"hqsjw4"`
		Sell4Volume int     `json:"hqssl4"`
		Sell5Price  float64 `json:"hqsjw5"`
		Sell5Volume int     `json:"hqssl5"`
	}
}

type bseChart struct {
	Data struct {
		Line []struct {
			Time  string  `json:"HQGXSJ"`
			Price float64 `json:"HQZJCJ"`
			Last  float64 `json:"HQZRSP"`
		}
	}
}

func (s *BSE) getRealtime() *BSE {
	s.Realtime.Index = "BSE"
	s.Realtime.Code = s.Code

	var res []bse
	resp, err := stock.Session.Get(api+s.Code, nil)
	if err != nil {
		log.Println("Failed to get bse realtime:", err)
		return s
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return s
	}
	if err := json.Unmarshal(trimCallback(resp.Bytes()), &res); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}
	if len(res) == 0 || len(res[0].Content) == 0 {
		log.Print("no result")
		return s
	}
	r := res[0].Content[0]

	s.Realtime.Name = r.Name
	s.Realtime.Now = r.Now
	s.Realtime.High = r.High
	s.Realtime.Low = r.Low
	s.Realtime.Open = r.Open
	s.Realtime.Last = r.Last
	s.Realtime.Change = r.Change
	s.Realtime.Percent = fmt.Sprintf("%g%%", r.Percent)
	s.Realtime.Update = time.Now().Format(time.RFC3339)

	s.Realtime.Buy5 = []stock.SellBuy{
		{Price: r.Buy1Price, Volume: r.Buy1Volume / 100},
		{Price: r.Buy2Price, Volume: r.Buy2Volume / 100},
		{Price: r.Buy3Price, Volume: r.Buy3Volume / 100},
		{Price: r.Buy4Price, Volume: r.Buy4Volume / 100},
		{Price: r.Buy5Price, Volume: r.Buy5Volume / 100},
	}
	s.Realtime.Sell5 = []stock.SellBuy{
		{Price: r.Sell1Price, Volume: r.Sell1Volume / 100},
		{Price: r.Sell2Price, Volume: r.Sell2Volume / 100},
		{Price: r.Sell3Price, Volume: r.Sell3Volume / 100},
		{Price: r.Sell4Price, Volume: r.Sell4Volume / 100},
		{Price: r.Sell5Price, Volume: r.Sell5Volume / 100},
	}

	if reflect.DeepEqual(s.Realtime.Sell5, []stock.SellBuy{{}, {}, {}, {}, {}}) &&
		reflect.DeepEqual(s.Realtime.Buy5, []stock.SellBuy{{}, {}, {}, {}, {}}) {
		s.Realtime.Buy5 = []stock.SellBuy{}
		s.Realtime.Sell5 = []stock.SellBuy{}
	}

	return s
}

func (s *BSE) getChart() *BSE {
	var r bseChart
	resp, err := stock.Session.Get(fmt.Sprintf(chartAPI, s.Code), nil)
	if err != nil {
		log.Println("Failed to get bse chart:", err)
		return s
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return s
	}
	if err := resp.JSON(&r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return s
	}
	if len(r.Data.Line) > 0 {
		s.Chart.Last = r.Data.Line[0].Last
	}

	sort.Slice(r.Data.Line, func(i, j int) bool { return r.Data.Line[i].Time < r.Data.Line[j].Time })

	for _, i := range r.Data.Line {
		if x := i.Time[:2] + ":" + i.Time[2:4]; x != "13:00" {
			s.Chart.Data = append(s.Chart.Data, stock.Point{X: x, Y: i.Price})
		}
	}

	return s
}

// GetRealtime gets the bse stock's realtime information.
func (s *BSE) GetRealtime() stock.Realtime {
	return s.getRealtime().Realtime
}

// GetChart gets the bse stock's chart
func (s *BSE) GetChart() stock.Chart {
	return s.getChart().Chart
}

func trimCallback(b []byte) []byte {
	b = bytes.TrimPrefix(b, []byte("null("))
	return bytes.TrimSuffix(b, []byte(")"))
}

// Suggests returns bse stock suggests according the keyword.
func Suggests(keyword string) (suggests []stock.Suggest) {
	resp, err := gohttp.Post(suggestAPI, nil, url.Values{"code": {keyword}})
	if err != nil {
		log.Println("Failed to get bse suggest:", err)
		return
	} else if resp.StatusCode != 200 {
		log.Println("Bad status code:", resp.StatusCode)
		return
	}
	var r []string
	if err := json.Unmarshal(trimCallback(resp.Bytes()), &r); err != nil {
		log.Println("Unmarshal json Error:", err)
		return
	}

	re := regexp.MustCompile(stock.BSEPattern)
	for _, i := range r {
		s := strings.Split(i, "##")
		if re.MatchString(s[0]) {
			suggests = append(suggests, stock.Suggest{
				Index: "bse",
				Code:  s[0],
				Name:  s[1],
			})
		}
	}

	return
}

func init() {
	stock.RegisterStock(
		"bse",
		stock.BSEPattern,
		func(code string) stock.Stock {
			return &BSE{Code: code}
		},
		Suggests,
	)
}
