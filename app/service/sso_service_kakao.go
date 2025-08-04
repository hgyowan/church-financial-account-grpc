package service

import (
	"context"
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain/kakao"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

type ssoKakaoService struct {
	s *service
}

func (s *ssoKakaoService) GetSSOUser(ctx context.Context, request user.GetSSOUserRequest) (*user.GetSSOUserResponse, error) {
	userInfo, err := s.s.client.GetUser(ctx, kakao.GetUserRequest{AccessToken: request.AccessToken})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if userInfo.KakaoAccount.IsEmailVerified == false || userInfo.KakaoAccount.Email == "" {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.InvalidSSOAccount)
	}

	return &user.GetSSOUserResponse{SSOUser: user.SSOUser{
		SSOUserID: fmt.Sprintf("%d", userInfo.ID),
		Nickname:  userInfo.KakaoAccount.Profile.Nickname,
		Email:     userInfo.KakaoAccount.Email,
	}}, nil
}

func (s *ssoKakaoService) IssueToken(ctx context.Context, request user.IssueTokenRequest) (*user.IssueTokenResponse, error) {
	kakaoToken, err := s.s.client.GetOauthToken(ctx, kakao.GetOauthTokenRequest{Code: request.Code})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &user.IssueTokenResponse{
		AccessToken:  kakaoToken.AccessToken,
		RefreshToken: kakaoToken.RefreshToken,
	}, nil
}

func NewSSOServiceKakao(s *service) *ssoKakaoService {
	return &ssoKakaoService{
		s: s,
	}
}
