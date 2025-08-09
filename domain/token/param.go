package token

import "time"

type GenerateClaimRequest struct {
	UserID             string        `json:"userId"`
	ExpireTimeDuration time.Duration `json:"expireTimeDuration"`
	Secret             string        `json:"secret"`
}

type IssueJWTTokenRequest struct {
	UserID string `json:"userId"`
}

type IssueJWTTokenResponse struct {
	JWTToken JWTToken `json:"jwtToken"`
}

type RefreshJWTTokenRequest struct {
	UserID       string `json:"userId"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJWTTokenResponse struct {
	JWTToken JWTToken `json:"jwtToken"`
}
