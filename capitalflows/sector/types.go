package sector

type Sectors interface {
	GetChart() []Chart
	GetTimeLine() []TimeLine
}

// TimeLine contains one day timeline data.
type TimeLine struct {
	Sector   string             `json:"sector"`
	TimeLine []map[string]int64 `json:"timeline"`
}

// Chart contains one day chart data.
type Chart struct {
	Sector string `json:"sector"`
	Chart  []XY   `json:"chart"`
}

type XY struct {
	X string `json:"x"`
	Y int64  `json:"y"`
}

// Chart2TimeLine convert Chart to TimeLine.
func Chart2TimeLine(chart Chart) TimeLine {
	var data []map[string]int64
	for _, point := range chart.Chart {
		data = append(data, map[string]int64{point.X: point.Y})
	}

	return TimeLine{Sector: chart.Sector, TimeLine: data}
}

// TimeLine2Chart convert TimeLine to Chart.
func TimeLine2Chart(timeline TimeLine) Chart {
	var xy []XY
	for _, i := range timeline.TimeLine {
		for k, v := range i {
			xy = append(xy, XY{k, v})
		}
	}

	return Chart{Sector: timeline.Sector, Chart: xy}
}
