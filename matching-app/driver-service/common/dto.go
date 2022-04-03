package common

type LocationDto struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DriverInRangeRequestDto struct {
	LocationDto
	RadiusInKm float64 `json:"radiusinkm"`
}
