package domain

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
)

type Service interface {
	user.UserService
	token.TokenService
	workspace.WorkspaceService
}
