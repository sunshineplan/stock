package sector

import "github.com/sunshineplan/database/mongodb"

var _ Sectors = charts{}

type charts []Chart

func (c charts) GetChart() []Chart {
	return c
}

func (c charts) GetTimeLine() (timelines []TimeLine) {
	for _, i := range c {
		timelines = append(timelines, Chart2TimeLine(i))
	}
	return
}

func GetSectors(date string, client mongodb.Client) (Sectors, error) {
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

	var charts charts
	for _, i := range data {
		charts = append(charts, Chart{i.ID, i.Chart})
	}
	return charts, nil
}
