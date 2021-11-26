package sector

import (
	"github.com/sunshineplan/database/mongodb"
)

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

func query(date string, xy bool, client mongodb.Client) (interface{}, error) {
	var data []struct {
		ID    string `json:"_id" bson:"_id"`
		Chart []XY
	}
	if err := client.Aggregate([]mongodb.M{
		{"$match": mongodb.M{"date": date}},
		{"$project": mongodb.M{"time": 1, "flows": mongodb.M{"$objectToArray": "$flows"}}},
		{"$unwind": "$flows"},
		{"$group": mongodb.M{
			"_id":   "$flows.k",
			"chart": mongodb.M{"$push": mongodb.M{"x": "$time", "y": "$flows.v"}},
		}},
		{"$sort": mongodb.M{"_id": 1}},
	}, &data); err != nil {
		return nil, err
	}

	var charts []Chart
	for _, i := range data {
		charts = append(charts, Chart{i.ID, i.Chart})
	}

	if xy {
		return charts, nil
	}

	var res []TimeLine
	for _, i := range charts {
		res = append(res, Chart2TimeLine(i))
	}

	return res, nil
}

// GetTimeLine gets all sectors timeline data of one day.
func GetTimeLine(date string, client mongodb.Client) ([]TimeLine, error) {
	res, err := query(date, false, client)
	if err != nil {
		return nil, err
	}

	return res.([]TimeLine), nil
}

// GetChart gets all sectors chart data of one day.
func GetChart(date string, client mongodb.Client) ([]Chart, error) {
	res, err := query(date, true, client)
	if err != nil {
		return nil, err
	}

	return res.([]Chart), nil
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
