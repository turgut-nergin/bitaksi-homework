package services

import (
	"encoding/json"
	"errors"

	"homework.matching-service/client"
	"homework.matching-service/common"
)

type MatchingService struct {
	client *client.DriverClient
}

func isValidLocationDto(locationData *common.MatchingRequestDto) error {
	latitude := locationData.Latitude
	longitude := locationData.Longitude

	if latitude > 180 || latitude < -180 {
		return errors.New("Not Valid latitude values")
	}

	if longitude > 90 || longitude < -90 {
		return errors.New("Not Valid longitude values")
	}
	return nil
}

func (driverClient *MatchingService) MathchingDriver(locationData *json.Decoder) (*common.DriverInRangeDto, error) {

	var locationDto common.MatchingRequestDto
	err := locationData.Decode(&locationDto)

	if err != nil {
		return nil, err
	}

	if err := isValidLocationDto(&locationDto); err != nil {
		return nil, err
	}

	driversDistance, err := driverClient.client.GetDriverDistance(&locationDto)
	if err != nil {
		return nil, err
	}

	nearestLocation, err := findNearestDriverLocaiton(driversDistance)
	if err != nil {
		return nil, err
	}
	return nearestLocation, nil

}

func findNearestDriverLocaiton(driversDistance []common.DriverInRangeDto) (*common.DriverInRangeDto, error) {
	nearestLocation := &driversDistance[0]
	for index, driverDistance := range driversDistance {
		if driverDistance.Distance < nearestLocation.Distance {
			nearestLocation = &driversDistance[index]
		}
	}

	return nearestLocation, nil
}

func New(driverClient *client.DriverClient) *MatchingService {
	driverService := MatchingService{driverClient}
	return &driverService
}
