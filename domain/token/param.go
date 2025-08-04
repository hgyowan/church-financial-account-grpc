package token

import "time"

type GenerateClaimRequest struct {
	UserID             string        `json:"user_id"`
	ExpireTimeDuration time.Duration `json:"expire_time_duration"`
	Secret             string        `json:"secret"`
}

type IssueJWTTokenRequest struct {
	UserID string `json:"user_id"`
}

type IssueJWTTokenResponse struct {
	JWTToken JWTToken `json:"jwt_token"`
}

type RefreshJWTTokenRequest struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshJWTTokenResponse struct {
	JWTToken JWTToken `json:"jwt_token"`
}
