package user

import "context"

type UserService interface {
	RegisterEmailUser(ctx context.Context, request RegisterEmailUserRequest) error
	SendVerifyEmail(ctx context.Context, request SendVerifyEmailRequest) error
	VerifyEmail(ctx context.Context, request VerifyEmailRequest) error
	LoginSSO(ctx context.Context, request LoginSSORequest) (*LoginSSOResponse, error)
	RegisterSSOUser(ctx context.Context, request RegisterSSOUserRequest) error
}

type SSOService interface {
	IssueToken(ctx context.Context, request IssueTokenRequest) (*IssueTokenResponse, error)
	GetSSOUser(ctx context.Context, request GetSSOUserRequest) (*GetSSOUserResponse, error)
}
