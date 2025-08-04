package external

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"resty.dev/v3"
)

type externalHttpClient struct {
	client *resty.Client
}

func (e *externalHttpClient) Client() *resty.Client {
	return e.client
}

func MustNewExternalHttpClient() domain.ExternalHttpClient {
	return &externalHttpClient{client: resty.New()}
}
