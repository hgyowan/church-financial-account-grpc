package service

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
)

func registerTokenService(s *service) {
	s.TokenService = &tokenService{s: s}
}

type tokenService struct {
	s *service
}

func (t *tokenService) IssueJWTToken(ctx context.Context, request token.IssueJWTTokenRequest) (*token.IssueJWTTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tokenService) RefreshJWTToken(ctx context.Context, request token.RefreshJWTTokenRequest) (*token.RefreshJWTTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
