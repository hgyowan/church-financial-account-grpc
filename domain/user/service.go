package user

import "context"

type UserService interface {
	SendVerifyEmail(ctx context.Context, request SendVerifyEmailRequest) error
	VerifyEmail(ctx context.Context, request VerifyEmailRequest) error
	RegisterEmailUser(ctx context.Context, request RegisterEmailUserRequest) error
	RegisterSSOUser(ctx context.Context, request RegisterSSOUserRequest) error
	LoginSSO(ctx context.Context, request LoginSSORequest) (*LoginSSOResponse, error)
	LoginEmail(ctx context.Context, request LoginEmailRequest) (*LoginEmailResponse, error)
	GetUser(ctx context.Context, request GetUserRequest) (*GetUserResponse, error)
	ListUserSimple(ctx context.Context, request ListUserSimpleRequest) (*ListUserSimpleResponse, error)
}

type SSOService interface {
	IssueToken(ctx context.Context, request IssueTokenRequest) (*IssueTokenResponse, error)
	GetSSOUser(ctx context.Context, request GetSSOUserRequest) (*GetSSOUserResponse, error)
}
