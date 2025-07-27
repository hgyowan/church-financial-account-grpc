package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

type userRepository struct {
	repository *repository
}

func registerUserRepository(r *repository) {
	r.UserRepository = &userRepository{repository: r}
}

func (u *userRepository) CreateUser(param *user.User) error {
	return pkgError.Wrap(u.repository.externalGormClient.DB().Create(param).Error)
}
