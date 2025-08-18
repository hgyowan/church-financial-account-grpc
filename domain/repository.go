package domain

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type Repository interface {
	user.UserRepository
	WithTransaction(fn func(txRepo Repository) error) error
}
