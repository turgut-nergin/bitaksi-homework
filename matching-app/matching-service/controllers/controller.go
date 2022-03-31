package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"homework.matching-service/services"
)

type MatchingController struct {
	matchService *services.MatchingService
}

func getErrorIfMethodNotAllowed(req *http.Request, method string) error {
	if req.Method != method {
		return errors.New("method not allowed")
	}
	return nil
}

func corsConfiguration(w http.ResponseWriter, mehtod string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Cache-Control")
	w.Header().Add("Access-Control-Allow-Methods", mehtod)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func (mC MatchingController) matchDriver(w http.ResponseWriter, req *http.Request) {
	corsConfiguration(w, "POST")
	err := getErrorIfMethodNotAllowed(req, http.MethodPost)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	decoder := json.NewDecoder(req.Body)
	nearestLocation, err := mC.matchService.MathchingDriver(decoder)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		return
	}
	nearestLocationEncode, err := json.Marshal(&nearestLocation)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}

	fmt.Fprintf(w, "%v\n nearest locaiton!", string(nearestLocationEncode))
}

func (mC MatchingController) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/matching", http.HandlerFunc(mC.matchDriver))
	return mux
}

func New(matchingService *services.MatchingService) *MatchingController {
	matchingController := MatchingController{matchingService}
	return &matchingController
}
