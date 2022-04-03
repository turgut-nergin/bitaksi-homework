package repository

import (
	"context"

	"homework.driver-service/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (dR *DriverRepository) BulkLocation(locations []common.LocationDto) error {
	var locationPointDto LocationPoint
	var result []interface{}
	collection := dR.db.Collection("drivers")
	for _, location := range locations {
		locationPointDto.Type = "Point"
		locationPointDto.Coordinates = []float64{location.Latitude, location.Longitude}
		result = append(result, geoJSON{locationPointDto})
	}
	_, err := collection.InsertMany(context.TODO(), result)
	if err != nil {
		return err
	}
	return nil
}

func (dR *DriverRepository) AddLocation(location LocationPoint) {
	collection := dR.db.Collection("drivers")
	geoJson := geoJSON{Location: location}
	collection.InsertOne(context.TODO(), geoJson)

}

func buildNearbyStage(location LocationPoint, radiusInKm float64) primitive.D {
	getNearbyStage := bson.D{{"$geoNear", bson.M{
		"near": bson.M{
			"type":        location.Type,
			"coordinates": location.Coordinates,
		},
		"maxDistance":        radiusInKm * 1000, //radius to meters
		"spherical":          true,
		"distanceField":      "distance",
		"distanceMultiplier": 0.001}}}
	return getNearbyStage
}

func (dR *DriverRepository) GetDriversDistance(radiusInKm float64, location LocationPoint) ([]DriverInRangeDto, error) {
	collection := dR.db.Collection("drivers")
	stages := mongo.Pipeline{}
	getNearbyStage := buildNearbyStage(location, radiusInKm)

	stages = append(stages, getNearbyStage)

	cursor, err := collection.Aggregate(context.TODO(), stages)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	var driverDistances []DriverInRangeDto
	cursor.All(context.TODO(), &driverDistances)

	return driverDistances, err

}

func New(db *mongo.Database) *DriverRepository {
	driverRepository := DriverRepository{db}
	return &driverRepository
}
