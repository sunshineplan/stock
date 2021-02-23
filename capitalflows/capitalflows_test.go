package capitalflows

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFlows(t *testing.T) {
	flows, err := Fetch()
	if err != nil {
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
