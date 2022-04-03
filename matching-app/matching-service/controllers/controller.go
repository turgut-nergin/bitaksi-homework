package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

var jwtKey = []byte("bi_taksi_api_key")

func corsConfiguration(w http.ResponseWriter, mehtod string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin,  Cache-Control, Authorization")
	w.Header().Add("Access-Control-Allow-Methods", mehtod)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

func (mC MatchingController) Signin(w http.ResponseWriter, r *http.Request) {
	corsConfiguration(w, "POST")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	if creds.Authenticated != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Fprintf(w, "%s\n", "sd")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	var loginConfig loginConfig
	loginConfig.Name = "token"
	loginConfig.Value = tokenString
	json.NewEncoder(w).Encode(loginConfig)

}

func (mC MatchingController) matchDriver(w http.ResponseWriter, req *http.Request) {
	corsConfiguration(w, "POST")

	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")

	if len(splitToken) != 2 {
		fmt.Fprintf(w, "Invalid Authorization")
		return
	}

	reqToken = strings.TrimSpace(splitToken[1])

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = getErrorIfMethodNotAllowed(req, http.MethodPost)

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
		return
	}
	fmt.Fprintf(w, "%v\n nearest locaiton!", string(nearestLocationEncode))
}

func (mC MatchingController) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/login", http.HandlerFunc(mC.Signin))
	mux.Handle("/matching", http.HandlerFunc(mC.matchDriver))
	return mux
}

func New(matchingService *services.MatchingService) *MatchingController {
	matchingController := MatchingController{matchingService}
	return &matchingController
}
