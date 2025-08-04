package token

import "context"

type TokenService interface {
	IssueJWTToken(ctx context.Context, request IssueJWTTokenRequest) (*IssueJWTTokenResponse, error)
	RefreshJWTToken(ctx context.Context, request RefreshJWTTokenRequest) (*RefreshJWTTokenResponse, error)
}
