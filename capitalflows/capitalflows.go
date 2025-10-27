package capitalflows

import (
	"fmt"

	"github.com/sunshineplan/stock"
)

const api = "http://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=500&fields=f14%2Cf62&fs=m%3A90%2Bt%3A2"

// CapitalFlows represents capital flows of all stock sectors.
type CapitalFlows map[string]int64

// Fetch fetchs capital flows.
func Fetch() (cf CapitalFlows, err error) {
	var res struct {
		Data struct {
			Diff map[string]struct {
				F14 string
				F62 float64
			}
			Total int
		}
	}
	resp, err := stock.Session.Get(api, nil)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code: %d", resp.StatusCode)
		return
	}
	if err = resp.JSON(&res); err != nil {
		return
	}

	cf = make(CapitalFlows)
	for _, v := range res.Data.Diff {
		cf[v.F14] = int64(v.F62)
	}

	return
}
