package user

import "context"

type UserService interface {
	CreateEmailUser(ctx context.Context, request CreateEmailUserRequest) error
	SendVerifyEmail(ctx context.Context, request SendVerifyEmailRequest) error
	VerifyEmail(ctx context.Context, request VerifyEmailRequest) error
}
