package controllers

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Authenticated bool `json:"authenticated"`
}

type Claims struct {
	jwt.StandardClaims
}

type loginConfig struct {
	Name  string
	Value string
}
