package sector

import (
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	chart := Chart{
		Sector: "Test",
		Chart: []struct {
			X string `json:"x"`
			Y int64  `json:"y"`
		}{
			{X: "a", Y: 1},
			{X: "b", Y: 2},
		},
	}
	timeline := TimeLine{
		Sector: "Test",
		TimeLine: []map[string]int64{
			{"a": 1},
			{"b": 2},
		},
	}

	if !reflect.DeepEqual(chart, TimeLine2Chart(timeline)) {
		t.Error("TimeLine2Chart got wrong result.")
	}
	if !reflect.DeepEqual(timeline, Chart2TimeLine(chart)) {
		t.Error("Chart2TimeLine got wrong result.")
	}
}
