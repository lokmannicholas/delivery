package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lokmannicholas/delivery/pkg/config"
)

type DistanceFounder interface {
	CountDistance(start, end []string) int
}

type DistanceFounderImp struct {
	apiKey  string
	rootUrl string
}

func GetDistanceFounder() DistanceFounder {
	return &DistanceFounderImp{
		rootUrl: "https://maps.googleapis.com/maps/api",
		apiKey:  config.Get().MapApiKey,
	}
}

type DistanceResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			DurationInTraffic struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration_in_traffic"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

func (d *DistanceFounderImp) CountDistance(start, end []string) int {

	url := fmt.Sprintf(`%s/distancematrix/json?origins=%s,%s&destinations=%s,%s&mode=driving&language=en&departure_time=now&key=%s`,
		d.rootUrl,
		start[0], start[1],
		end[0], end[1],
		d.apiKey)
	res, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	r := new(DistanceResponse)
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range r.Rows {
		for _, el := range row.Elements {
			if el.Status == "OK" {
				return el.Distance.Value
			}
		}
	}
	return 0
}
