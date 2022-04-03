package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverRepository struct { //It must be move, driver repository struct is not dto
	db *mongo.Database
}

type LocationPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type geoJSON struct {
	Location LocationPoint
}

type DriverInRangeDto struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Distance float64            `bson:"distance, omitempty"`
	Location LocationPoint      `bson:"location, omitempty"`
}
