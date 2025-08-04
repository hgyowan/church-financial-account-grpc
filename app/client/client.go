package client

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/kakao"
)

type client struct {
	kakao.KakaoClient

	externalHttpClient domain.ExternalHttpClient
}

func NewClient(externalHttpClient domain.ExternalHttpClient) domain.Client {
	c := &client{
		externalHttpClient: externalHttpClient,
	}

	c.register()

	return c
}

func (c *client) register() {
	registerKakaoClient(c)
}
