package capitalflows

import (
	"log"
	"testing"

	"github.com/sunshineplan/utils/retry"
)

func TestFlows(t *testing.T) {
	var flows CapitalFlows
	if err := retry.Do(func() (err error) {
		flows, err = Fetch()
		return
	}, 5, 20); err != nil {
		t.Fatal(err)
	}
	log.Println(len(flows), flows)
}
