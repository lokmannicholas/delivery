package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/lokmannicholas/delivery/pkg/config"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DistanceFounder interface {
	CountDistance(start, end []string) int
}

type DistanceFounderImp struct {
	apiKey  string
	rootUrl string
	client  HttpClient
}

func GetDistanceFounder() DistanceFounder {
	return &DistanceFounderImp{
		rootUrl: "https://maps.googleapis.com/maps/api",
		apiKey:  config.Get().MapApiKey,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type DistanceResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []Row    `json:"rows"`
	Status               string   `json:"status"`
}
type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance          ElementsInfo `json:"distance"`
	Duration          ElementsInfo `json:"duration"`
	DurationInTraffic ElementsInfo `json:"duration_in_traffic"`
	Status            string       `json:"status"`
}
type ElementsInfo struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

func (d *DistanceFounderImp) CountDistance(start, end []string) int {

	url := fmt.Sprintf(`%s/distancematrix/json?origins=%s,%s&destinations=%s,%s&mode=driving&language=en&departure_time=now&key=%s`,
		d.rootUrl,
		start[0], start[1],
		end[0], end[1],
		d.apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := d.client.Do(req)

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
