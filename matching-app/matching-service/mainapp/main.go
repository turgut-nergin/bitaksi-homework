package main

import (
	"log"
	"net/http"

	"homework.matching-service/client"
	"homework.matching-service/controllers"
	"homework.matching-service/services"
)

func main() {
	driverClient := client.New("BITAKSI_API_KEY")
	matchingService := services.New(driverClient)
	matchingController := controllers.New(matchingService)
	log.Println("- Matching Service 8081 -")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", matchingController.GetRoutes()))
}
