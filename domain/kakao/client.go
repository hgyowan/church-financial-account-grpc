package kakao

import "context"

type KakaoClient interface {
	GetOauthToken(ctx context.Context, request GetOauthTokenRequest) (*GetOAuthTokenResponse, error)
	GetUser(ctx context.Context, request GetUserRequest) (*GetUserResponse, error)
}
