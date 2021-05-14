package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID       uint64
	Username string
	Role     int32
	jwt.StandardClaims
}
