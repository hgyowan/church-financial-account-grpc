package client

import (
	"context"
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain/kakao"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"io"
)

const (
	kakaoAuthHost = "https://kauth.kakao.com"
	kakaoAPIHost  = "https://kapi.kakao.com"
)

func registerKakaoClient(c *client) {
	c.KakaoClient = &kakaoClient{
		client: c,
	}
}

type kakaoClient struct {
	client *client
}

func (k *kakaoClient) GetUser(ctx context.Context, request kakao.GetUserRequest) (*kakao.GetUserResponse, error) {
	var userResponse *kakao.GetUserResponse

	res, err := k.client.externalHttpClient.Client().R().
		SetResult(&userResponse).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=utf-8").
		SetAuthToken(request.AccessToken).
		Get(kakaoAPIHost + "/v2/user/me")
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	data, _ := io.ReadAll(res.Body)
	fmt.Println(string(data))
	if res.StatusCode() != 200 {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.Kakao)
	}

	return userResponse, nil
}

func (k *kakaoClient) GetOauthToken(ctx context.Context, request kakao.GetOauthTokenRequest) (*kakao.GetOAuthTokenResponse, error) {
	var oauthTokenResponse *kakao.GetOAuthTokenResponse
	res, err := k.client.externalHttpClient.Client().R().
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
		Post(kakaoAuthHost + "/oauth/token")
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if res.StatusCode() != 200 {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.Kakao)
	}

	return oauthTokenResponse, nil
}
