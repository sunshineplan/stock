package stock

import (
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sunshineplan/gohttp"
	"github.com/sunshineplan/utils"
)

const (
	SSEPattern  = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`
	SZSEPattern = `(00[0-3]|159|30[0-1]|399)\d{3}`
)

var Session = gohttp.NewSession()

// SetTimeout sets Session's timeout.
func SetTimeout(duration int) {
	if duration <= 0 {
		panic("duration must greater than 0")
	}
	Session.SetTimeout(time.Duration(duration) * time.Second)
}

// Stock represents a stock.
type Stock interface {
	GetRealtime() Realtime
	GetChart() Chart
}

// A format holds an stock format's index, pattern and how to decode it.
type format struct {
	index    string
	pattern  string
	init     func(string) Stock
	suggests func(string) []Suggest
}

// Formats is the list of registered formats.
var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterStock registers an stock format for use by Decode.
func RegisterStock(index, pattern string, init func(string) Stock, suggests func(string) []Suggest) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{index, pattern, init, suggests}))
	formatsMu.Unlock()
}

// Init initializes a stock.
func Init(index, code string) Stock {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		if strings.ToLower(index) == f.index {
			re := regexp.MustCompile(f.pattern)
			if re.MatchString(code) {
				return f.init(code)
			}
		}
	}
	return nil
}

// Realtime is a stock's realtime information.
type Realtime struct {
	Index   string    `json:"index"`
	Code    string    `json:"code"`
	Name    string    `json:"name"`
	Now     float64   `json:"now"`
	Change  float64   `json:"change"`
	Percent string    `json:"percent"`
	Sell5   []SellBuy `json:"sell5"`
	Buy5    []SellBuy `json:"buy5"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Open    float64   `json:"open"`
	Last    float64   `json:"last"`
	Update  string    `json:"update"`
}

// SellBuy represents stock's realtime sell buy information.
type SellBuy struct {
	Price  float64
	Volume int
}

// Realtimes returns stocks's realtime.
func Realtimes(s []Stock) []Realtime {
	r := make([]Realtime, len(s))
	var wg sync.WaitGroup
	for i, v := range s {
		wg.Add(1)
		go func(i int, s Stock) {
			defer wg.Done()
			if s != nil {
				r[i] = s.GetRealtime()
			}
		}(i, v)
	}
	wg.Wait()
	return r
}

// Chart is a stock's chart data.
type Chart struct {
	Last float64 `json:"last"`
	Data []Point `json:"chart"`
}

// Point represents stock's chart point information.
type Point struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

// Suggest represents stock suggest.
type Suggest struct {
	Index string
	Code  string
	Name  string
	Type  string
}

// Suggests suggests stocks according keyword.
func Suggests(keyword string) []Suggest {
	formats, _ := atomicFormats.Load().([]format)
	suggests := []Suggest{}
	for _, f := range formats {
		suggests = append(suggests, f.suggests(keyword)...)
	}
	return utils.Deduplicate(suggests).([]Suggest)
}
