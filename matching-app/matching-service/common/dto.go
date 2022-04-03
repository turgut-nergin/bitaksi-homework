package common

type MatchingRequestDto struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	RadiusInKm float64 `json:"radiusinkm"`
}

type LocationPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type DriverInRangeDto struct {
	ID       string        `bson:"_id,omitempty"`
	Distance float64       `bson:"distance, omitempty"`
	Location LocationPoint `bson:"location, omitempty"`
}
