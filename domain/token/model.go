package token

import "github.com/dgrijalva/jwt-go"

type JWTCustomClaims struct {
	UserID uint
	jwt.StandardClaims
}
