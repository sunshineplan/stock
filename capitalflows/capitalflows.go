package capitalflows

import "github.com/sunshineplan/gohttp"

const api = "http://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=500&fields=f14%2Cf62&fs=m%3A90%2Bt%3A2"

// CapitalFlows represents capital flows of a stock sector.
type CapitalFlows map[string]int64

// Fetch fetchs capital flows of all stock sectors.
func Fetch() (cf []CapitalFlows, err error) {
	var res struct {
		Data struct {
			Diff map[string]struct {
				F14 string
				F62 float64
			}
			Total int
		}
	}
	if err = gohttp.Get(api, nil).JSON(&res); err != nil {
		return
	}

	for _, v := range res.Data.Diff {
		cf = append(cf, CapitalFlows{v.F14: int64(v.F62)})
	}

	return
}
