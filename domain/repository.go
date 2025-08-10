package domain

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
)

type Repository interface {
	user.UserRepository
	workspace.WorkspaceRepository
	WithTransaction(fn func(txRepo Repository) error) error
}
