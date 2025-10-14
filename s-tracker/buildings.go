package s_tracker

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/pterm/pterm"
	"math"
	"strconv"
	"strings"
)

const (
	buildingsApiPath  = "api/building"
	buildingsInfoPath = "studiebolig/building"
	countPerPage      = 10
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

type BuildingFull struct {
	BasicBuildingInfo Building

	BestWaitingListCategory rune
	BuildYear               int

	// TODO cus i dont care about these tbh
	HasSharedKitchen bool
	HasOwnKitchen    bool
	HasOwnToilet     bool
	HasOwnBath       bool
	HasElevator      bool
	HasCommonRoom    bool
	HasLaundry       bool
	IsSmokeFree      bool
}

func (c Client) GetUserBuildings() ([]Building, error) {
	var err error
	var buildings []Building
	var count int

	progressBar, _ := pterm.DefaultProgressbar.WithTotal(0).Start()

	c.Collector.OnResponse(func(r *colly.Response) {
		var buildingList BuildingList
		err = json.Unmarshal(r.Body, &buildingList)
		count = buildingList.Count
		buildings = append(buildings, buildingList.Results...)

		progressBar.Total = count
		progressBar.Add(len(buildingList.Results))
	})

	_ = c.Collector.Visit(fmt.Sprintf(
		"%s/%s/?has_application_for=%s&parent=1&page=%d",
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
			"%s/%s/?has_application_for=%s&parent=1&page=%d",
			SBaseUrl,
			buildingsApiPath,
			c.applicantId,
			i,
		))
	}

	pterm.Println()
	return buildings, nil
}

func (c Client) GetUserBuildingsFull(buildingList []Building) ([]BuildingFull, error) {
	currentBuildingIndex := 0
	buildingsFull := make([]BuildingFull, len(buildingList))

	pterm.Info.Println("Getting full building information...")
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(buildingsFull)).Start()

	// get all the shit (why couldnt u have given me an api im gonna kms)
	c.Collector.OnHTML(".form", func(e *colly.HTMLElement) {
		currentBuilding := &buildingsFull[currentBuildingIndex]
		currentBuilding.BasicBuildingInfo = buildingList[currentBuildingIndex]

		// waiting list category stuff
		bestWaitingListCategory := 'G'
		e.ForEach(".waiting-list-category", func(_ int, child *colly.HTMLElement) {
			category := []rune(strings.TrimSpace(child.Text))[0]
			if category < bestWaitingListCategory {
				bestWaitingListCategory = category
			}
		})
		currentBuilding.BestWaitingListCategory = bestWaitingListCategory

		// build year
		e.DOM.Find("dt").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Build year" {
				buildYear, err := strconv.Atoi(s.Next().Text())
				if err != nil {
					pterm.Warning.Printf("Unable to parse build year for building: %d", currentBuilding.BasicBuildingInfo.Id)
				}

				currentBuilding.BuildYear = buildYear
			}
		})

		// TODO: do the rest if i care some day (i probably wont lol)
		progressBar.Increment()
		currentBuildingIndex++
	})

	for _, building := range buildingList {
		_ = c.Collector.Visit(fmt.Sprintf("%s/%s/%d", SBaseUrl, buildingsInfoPath, building.Id))
	}

	c.Collector.OnHTMLDetach(".container")

	return buildingsFull, nil
}
