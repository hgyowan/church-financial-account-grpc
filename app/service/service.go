package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type service struct {
	user.UserService
	token.TokenService

	repo   domain.Repository
	client domain.Client

	externalRedisClient  domain.ExternalRedisClient
	externalMailSender   domain.ExternalMailSender
	externalSearchEngine domain.ExternalSearchEngine
	validator            domain.ExternalValidator
	grpcClient           domain.ExternalGRPCClient
}

func NewService(repo domain.Repository, client domain.Client, externalRedisClient domain.ExternalRedisClient, externalMailSender domain.ExternalMailSender, externalSearchEngine domain.ExternalSearchEngine, validator domain.ExternalValidator, grpcClient domain.ExternalGRPCClient) domain.Service {
	s := &service{
		repo:                 repo,
		client:               client,
		externalRedisClient:  externalRedisClient,
		externalMailSender:   externalMailSender,
		externalSearchEngine: externalSearchEngine,
		validator:            validator,
		grpcClient:           grpcClient,
	}
	s.register()
	return s
}

func (s *service) register() {
	registerUserService(s)
	registerTokenService(s)
}
