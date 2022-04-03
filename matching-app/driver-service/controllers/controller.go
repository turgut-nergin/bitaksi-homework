package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"homework.driver-service/services"
)

type DriverController struct {
	driverService *services.DriverService
}

func corsConfiguration(w http.ResponseWriter, mehtod string) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Cache-Control")
	header.Add("Access-Control-Allow-Methods", mehtod)
	header.Add("Access-Control-Allow-Credentials", "true")
}

func (dC DriverController) insertOne(w http.ResponseWriter, req *http.Request) {
	corsConfiguration(w, "POST")
	err := getErrorIfMethodNotAllowed(req, http.MethodPost)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = dC.driverService.InsertOne(decoder)
	if err != nil {
		fmt.Fprintf(w, "%s\n ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (dC DriverController) insertMany(w http.ResponseWriter, req *http.Request) {
	corsConfiguration(w, "POST")
	err := getErrorIfMethodNotAllowed(req, http.MethodPost)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	fourMbInByte := 1 << 22
	req.ParseMultipartForm(int64(fourMbInByte))
	file, _, err := req.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	err = dC.driverService.BulkDriver(&file)

	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (dC DriverController) getDriversInRange(w http.ResponseWriter, req *http.Request) {

	err := getErrorIfMethodNotAllowed(req, http.MethodPost)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	apiKey := req.Header.Get("apiKey")
	if apiKey != "bi_taksi_api_key" {
		fmt.Fprintf(w, "%s\n", "Permission Denied!")
		return
	}
	decoder := json.NewDecoder(req.Body)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	driversDistanceResult, err := dC.driverService.GetDrivers(decoder)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}

	json.NewEncoder(w).Encode(driversDistanceResult)
	return
}

func getErrorIfMethodNotAllowed(req *http.Request, method string) error {
	if req.Method != method {
		return errors.New("method not allowed")
	}
	return nil
}

func (dC DriverController) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/driver/insert", http.HandlerFunc(dC.insertOne))
	mux.Handle("/driver/insert/bulk", http.HandlerFunc(dC.insertMany))
	mux.Handle("/driver/inrange/sorted", http.HandlerFunc(dC.getDriversInRange))
	return mux
}

func New(driverService *services.DriverService) *DriverController {
	driverController := DriverController{driverService}
	return &driverController
}
