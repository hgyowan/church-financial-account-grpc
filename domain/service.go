package domain

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type Service interface {
	user.UserService
	token.TokenService
}
