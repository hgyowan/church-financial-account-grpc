package service

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type service struct {
	user.UserService

	repo domain.Repository

	externalRedisClient domain.ExternalRedisClient
}

func NewService(ctx context.Context, repo domain.Repository, externalRedisClient domain.ExternalRedisClient) domain.Service {
	s := &service{
		repo:                repo,
		externalRedisClient: externalRedisClient,
	}
	s.register(ctx)
	return s
}

func (s *service) register(ctx context.Context) {
	registerUserService(ctx, s)
}
