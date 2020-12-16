package stock

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicate(t *testing.T) {
	tc := []Suggest{
		{Index: "SSE", Code: "600309", Name: "万华化学", Type: "GP-A"},
		{Index: "SSE", Code: "600309", Name: "万华化学", Type: "GP-A"},
		{Index: "SSE", Code: "600519", Name: "贵州茅台", Type: "GP-A"},
	}
	expect := []Suggest{
		{Index: "SSE", Code: "600309", Name: "万华化学", Type: "GP-A"},
		{Index: "SSE", Code: "600519", Name: "贵州茅台", Type: "GP-A"},
	}
	unique := removeDuplicate(tc)
	if !reflect.DeepEqual(unique, expect) {
		t.Errorf("expected %v; got %v", expect, unique)
	}
}
