package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type repository struct {
	user.UserRepository
	externalGormClient domain.ExternalDBClient
}

func NewRepository(externalGormClient domain.ExternalDBClient) domain.Repository {
	r := &repository{
		externalGormClient: externalGormClient,
	}
	r.register()

	return r
}

func (r *repository) register() {
	registerUserRepository(r)
}
