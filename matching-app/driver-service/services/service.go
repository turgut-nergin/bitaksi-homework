package services

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"homework.driver-service/common"
	"homework.driver-service/repository"
)

type DriverService struct {
	driverRepository *repository.DriverRepository
}

type cvsFileLocationColumnInfo struct {
	LatitudeIndex  int
	LongitudeIndex int
}

func (d *DriverService) InsertOne(locationData *json.Decoder) error {
	var locationDto common.LocationDto
	err := locationData.Decode(&locationDto)
	if err != nil {
		return err
	}

	if locationDto.Latitude > 180 || locationDto.Latitude < -180 {
		return errors.New("Not Valid latitude values!")
	}

	if locationDto.Longitude > 90 || locationDto.Longitude < -90 {
		return errors.New("Not Valid latitude values!")
	}

	var locationPoint repository.LocationPoint
	locationPoint.Type = "Point"
	locationPoint.Coordinates = []float64{locationDto.Latitude, locationDto.Longitude}

	d.driverRepository.AddLocation(locationPoint)
	return nil
}

func (d *DriverService) GetDrivers(locaitonData *json.Decoder) ([]repository.DriverInRangeDto, error) {
	var locationDto common.DriverInRangeRequestDto
	err := locaitonData.Decode(&locationDto)
	if err != nil {
		return nil, err
	}

	var locationPoint repository.LocationPoint
	locationPoint.Type = "Point"
	locationPoint.Coordinates = []float64{locationDto.Latitude, locationDto.Longitude}
	driversDistance, err := d.driverRepository.GetDriversDistance(locationDto.RadiusInKm, locationPoint)
	if err != nil {
		return nil, err
	}
	return driversDistance, nil
}

func isLatLongValid(latLong []string, locationRowIndex *cvsFileLocationColumnInfo) (bool, []float64) {
	latitude, err := strconv.ParseFloat(latLong[locationRowIndex.LatitudeIndex], 64)
	if err != nil {
		log.Fatal("Not Valid latitude values")
		return false, nil
	}

	if latitude > 180 || latitude < -180 {
		log.Fatal("Not Valid latitude values")
		return false, nil
	}

	longitude, err := strconv.ParseFloat(latLong[locationRowIndex.LongitudeIndex], 64)

	if err != nil {
		log.Fatal("Not Valid latitude values")
		return false, nil
	}

	if longitude > 90 || longitude < -90 {
		log.Fatal("Not Valid latitude values")
		return false, nil
	}

	return true, []float64{latitude, longitude}

}

func getLatAndLongRowIndex(heads []string) (*cvsFileLocationColumnInfo, error) {

	var locationRowIndex cvsFileLocationColumnInfo

	if isLatitude := strings.ContainsAny(heads[0], "lat"); isLatitude {
		locationRowIndex.LatitudeIndex = 0
		locationRowIndex.LongitudeIndex = 1
		return &locationRowIndex, nil
	}

	if isLongitude := strings.ContainsAny(heads[0], "long"); isLongitude {
		locationRowIndex.LatitudeIndex = 1
		locationRowIndex.LongitudeIndex = 0
		return &locationRowIndex, nil
	}
	return nil, errors.New("The file can not be parsed!")
}

func parseCvsToLocationDto(locationLinesCsvr [][]string, locationRowIndex *cvsFileLocationColumnInfo) []common.LocationDto {
	var locations []common.LocationDto
	var location common.LocationDto
	for _, line := range locationLinesCsvr {
		isValid, latLong := isLatLongValid(line, locationRowIndex)
		if isValid {
			location.Latitude = latLong[locationRowIndex.LatitudeIndex]
			location.Longitude = latLong[locationRowIndex.LongitudeIndex]
			locations = append(locations, location)
		}

	}
	return locations
}

func (d *DriverService) BulkDriver(file *multipart.File) error {
	locationLinesCsvr, err := csv.NewReader(*file).ReadAll()

	if err != nil {
		return nil
	}

	heads := locationLinesCsvr[0]
	locationRowIndex, err := getLatAndLongRowIndex(heads)

	if err != nil {
		return err
	}

	locationLinesCsvrWithoutHead := locationLinesCsvr[1:]
	locations := parseCvsToLocationDto(locationLinesCsvrWithoutHead, locationRowIndex)

	err = d.driverRepository.BulkLocation(locations)
	if err != nil {
		return err
	}

	return nil
}

func New(driverRepository *repository.DriverRepository) *DriverService {
	driverService := DriverService{driverRepository}
	return &driverService
}
