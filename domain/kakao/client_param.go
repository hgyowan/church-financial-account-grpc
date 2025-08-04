package kakao

type GetOauthTokenRequest struct {
	Code string `json:"code"`
}

type GetOAuthTokenResponse struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	IDToken               string `json:"id_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type GetUserRequest struct {
	AccessToken string
}

type GetUserResponse struct {
	ID           uint64       `json:"id"`
	KakaoAccount KakaoAccount `json:"kakao_account"`
}
