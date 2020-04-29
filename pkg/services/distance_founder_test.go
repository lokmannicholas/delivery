package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"

	"github.com/lokmannicholas/delivery/pkg"
	"github.com/lokmannicholas/delivery/pkg/services/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type DistanceFounderTestSuite struct {
	suite.Suite
}

func TestDistanceFounderTestSuite(t *testing.T) {
	suite.Run(t, new(DistanceFounderTestSuite))
}

func (suite *DistanceFounderTestSuite) TestGetDistanceFounder() {
	t := suite.T()
	client := &mocks.HttpClient{}
	t.Run("DistanceFounderImp calls DistanceFounder", func(t *testing.T) {
		mgr := &DistanceFounderImp{
			client: client,
		}
		assert.NotNil(t, mgr, "no error")
	})
}
func (suite *DistanceFounderTestSuite) TestDistanceFounderImp_CountDistance() {
	t := suite.T()
	start := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}
	end := []string{pkg.FloatToString(-90 + rand.Float64()*(90-(-90))), pkg.FloatToString(-180 + rand.Float64()*(180-(-180)))}
	dR := &DistanceResponse{
		Rows: []Row{
			{
				Elements: []Element{
					{
						Status: "OK",
						Distance: ElementsInfo{
							Text:  "5.3 km",
							Value: 5308,
						},
						Duration: ElementsInfo{
							Text:  "18 mins",
							Value: 1085,
						},
						DurationInTraffic: ElementsInfo{
							Text:  "13 mins",
							Value: 751,
						},
					},
				},
			},
		},
	}

	var httpClientMock = &mocks.HttpClient{}
	mgr := &DistanceFounderImp{
		rootUrl: "https://maps.googleapis.com/maps/api",
		client:  httpClientMock,
		apiKey:  "",
	}
	url := fmt.Sprintf(`%s/distancematrix/json?origins=%s,%s&destinations=%s,%s&mode=driving&language=en&departure_time=now&key=%s`,
		mgr.rootUrl,
		start[0], start[1],
		end[0], end[1],
		mgr.apiKey)
	req, _ := http.NewRequest("GET", url, nil)

	t.Run("DistanceFounderImp calls DistanceFounder", func(t *testing.T) {
		data, _ := json.Marshal(dR)
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		}
		httpClientMock.On("Do", req).Return(response, nil)
		distance := mgr.CountDistance(start, end)
		assert.NotEqual(t, 0, distance)
		for _, row := range dR.Rows {
			for _, el := range row.Elements {
				if el.Status == "OK" {
					assert.Equal(t, el.Distance.Value, distance)
				}
			}
		}
		httpClientMock.AssertCalled(t, "Do", req)
	})

	t.Run("Should return errors when http client failed", func(t *testing.T) {
		httpClientMock.On("Do", req).Return(nil, nil)
		distance := mgr.CountDistance(start, end)
		assert.Equal(t, 0, distance)
		httpClientMock.AssertCalled(t, "Do", req)
	})
	t.Run("Should return errors when coordinate incorrect", func(t *testing.T) {
		dR := &DistanceResponse{
			Rows: []Row{
				{
					Elements: []Element{
						{
							Status: "NOT_FOUND",
						},
					},
				},
			},
		}
		data, _ := json.Marshal(dR)
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(data)),
		}
		httpClientMock.On("Do", req).Return(response, nil)
		distance := mgr.CountDistance(start, end)
		assert.Equal(t, 0, distance)
		for _, row := range dR.Rows {
			for _, el := range row.Elements {
				assert.Contains(t, el.Status, "NOT_FOUND")
			}
		}
		httpClientMock.AssertCalled(t, "Do", req)
	})
}
