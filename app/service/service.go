package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type service struct {
	user.UserService

	repo domain.Repository

	externalRedisClient domain.ExternalRedisClient
}

func NewService(repo domain.Repository, externalRedisClient domain.ExternalRedisClient) domain.Service {
	s := &service{
		repo:                repo,
		externalRedisClient: externalRedisClient,
	}
	s.register()
	return s
}

func (s *service) register() {
	registerUserService(s)
}
