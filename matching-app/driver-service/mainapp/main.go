package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"homework.driver-service/controllers"
	"homework.driver-service/repository"
	"homework.driver-service/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func getDatabase() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
	database := client.Database("driver")
	models := mongo.IndexModel{Keys: bsonx.Doc{{Key: "location", Value: bsonx.String("2dsphere")}}}
	opts := options.CreateIndexes().SetMaxTime(20 * time.Second)
	database.Collection("drivers").Indexes().CreateOne(context.Background(), models, opts)
	return database
}

func main() {
	database := getDatabase()
	driverRepository := repository.New(database)
	driverService := services.New(driverRepository)
	driverController := controllers.New(driverService)
	log.Fatal(http.ListenAndServe("driver-service:8000", driverController.GetRoutes()))
}
