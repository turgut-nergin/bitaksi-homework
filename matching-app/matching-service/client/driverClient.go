package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"homework.matching-service/common"
)

var baseURL = url.URL{
	Scheme: "http",
	Host:   "driver-service:8000",
	Path:   "/driver/",
}

type DriverClient struct {
	client *http.Client
	apiKey string
}

func (client *DriverClient) GetDriverDistance(locationData *common.MatchingRequestDto) ([]common.DriverInRangeDto, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(locationData)
	if err != nil {
		return nil, err
	}
	endpt := baseURL.ResolveReference(
		&url.URL{Path: "inrange/sorted"})
	req, err := http.NewRequest("POST", endpt.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("apikey", client.apiKey)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var driverDistances []common.DriverInRangeDto

	if err := decoder.Decode(&driverDistances); err != nil {
		fmt.Println(decoder)
		return nil, err
	}

	if len(driverDistances) == 0 {
		return nil, errors.New("404 - Not Found")
	}

	return driverDistances, nil

}

func New(apiKey string) *DriverClient {
	c := &http.Client{Timeout: time.Minute}
	return &DriverClient{
		client: c,
		apiKey: apiKey,
	}
}
