package service

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"resty.dev/v3"
)

const (
	kakaoHost = "https://kauth.kakao.com"
)

type ssoKakaoService struct {
	s *service
}

func (s *ssoKakaoService) GetSSOUser(ctx context.Context, request user.GetSSOUserRequest) (*user.GetSSOUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ssoKakaoService) IssueToken(ctx context.Context, request user.IssueTokenRequest) (*user.IssueTokenResponse, error) {
	client := resty.New()

	type OAuthTokenResponse struct {
		TokenType             string `json:"token_type"`               // 예: "bearer"
		AccessToken           string `json:"access_token"`             // 사용자 액세스 토큰
		IDToken               string `json:"id_token,omitempty"`       // OpenID Connect ID 토큰 (조건부)
		ExpiresIn             int    `json:"expires_in"`               // access_token 과 id_token 만료 시간 (초)
		RefreshToken          string `json:"refresh_token"`            // 리프레시 토큰
		RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"` // 리프레시 토큰 만료 시간 (초)
		Scope                 string `json:"scope,omitempty"`          // 예: "profile email openid" (공백 구분)
	}

	oauthTokenResponse := &OAuthTokenResponse{}
	_, err := client.R().
		SetFormData(map[string]string{
			"grant_type":    envs.KakaoGranTType,
			"client_id":     envs.KakaoClientID,
			"redirect_uri":  envs.KakaoRedirectURI,
			"code":          request.Code,
			"client_secret": envs.KakaoClientSecret,
		},
		).
		SetResult(&oauthTokenResponse).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=utf-8").
		Post(kakaoHost + "/oauth/token")
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	return &user.IssueTokenResponse{
		AccessToken:  oauthTokenResponse.AccessToken,
		RefreshToken: oauthTokenResponse.RefreshToken,
	}, nil
}

func NewSSOServiceKakao(s *service) *ssoKakaoService {
	return &ssoKakaoService{
		s: s,
	}
}
