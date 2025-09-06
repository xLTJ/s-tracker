package s_tracker

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"math"
)

const (
	buildingsApiPath = "/api/building"
	countPerPage     = 10
)

type BuildingList struct {
	Count   int        `json:"count"`
	Next    string     `json:"next"`
	Results []Building `json:"results"`
}

type Building struct {
	Id           int    `json:"pk"`
	Name         string `json:"name"`
	Address      string `json:"desc_address"`
	Municipality string `json:"municipality"`
}

func (c Client) GetUserBuildings() ([]Building, error) {
	var err error
	var buildings []Building
	var count int
	c.Collector.OnResponse(func(r *colly.Response) {
		var buildingList BuildingList
		err = json.Unmarshal(r.Body, &buildingList)
		count = buildingList.Count
		buildings = append(buildings, buildingList.Results...)
	})

	_ = c.Collector.Visit(fmt.Sprintf(
		"%s%s/?has_application_for=%s&parent=1&page=%d",
		SBaseUrl,
		buildingsApiPath,
		c.applicantId,
		1,
	))

	if err != nil {
		return []Building{}, fmt.Errorf("error decoding response: %v", err)
	}

	totalPages := math.Ceil(float64(count) / float64(countPerPage))

	for i := 2; i <= int(totalPages); i++ {
		_ = c.Collector.Visit(fmt.Sprintf(
			"%s%s/?has_application_for=%s&parent=1&page=%d",
			SBaseUrl,
			buildingsApiPath,
			c.applicantId,
			i,
		))
	}

	return buildings, nil
}
