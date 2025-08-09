package token

import "github.com/dgrijalva/jwt-go"

type JWTCustomClaims struct {
	UserID string
	jwt.StandardClaims
}

type JWTToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
