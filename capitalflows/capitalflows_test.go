package capitalflows

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sunshineplan/utils/retry"
)

func TestFlows(t *testing.T) {
	var flows CapitalFlows
	if err := retry.Do(func() (e error) {
		flows, e = Fetch()
		return
	}, 5, 20); err != nil {
		t.Fatal(err)
	}

	v := reflect.ValueOf(flows)
	T := v.Type()
	fmt.Println("Fields:", v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Int() == 0 {
			fmt.Println(T.Field(i).Name, 0)
		}
	}
}
