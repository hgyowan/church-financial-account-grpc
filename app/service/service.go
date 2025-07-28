package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type service struct {
	user.UserService

	repo domain.Repository

	externalRedisClient domain.ExternalRedisClient
	validator           domain.ExternalValidator
}

func NewService(repo domain.Repository, externalRedisClient domain.ExternalRedisClient, validator domain.ExternalValidator) domain.Service {
	s := &service{
		repo:                repo,
		externalRedisClient: externalRedisClient,
		validator:           validator,
	}
	s.register()
	return s
}

func (s *service) register() {
	registerUserService(s)
}
