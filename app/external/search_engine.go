package external

import (
	"context"
	"errors"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgZincSearch "github.com/hgyowan/go-pkg-library/zincsearch"
	pkgZincSearchModel "github.com/hgyowan/go-pkg-library/zincsearch/model"
)

type externalSearchEngineClient struct {
	zincSearch pkgZincSearch.ZincSearch
}

func (e *externalSearchEngineClient) ZincSearch() pkgZincSearch.ZincSearch {
	return e.zincSearch
}

func MustNewSearchEngine(ctx context.Context) domain.ExternalSearchEngine {
	zincSearchCli := pkgZincSearch.MustNewZincSearch(ctx, &pkgZincSearch.ZinSearchConfig{
		Host:     envs.ZincSearchHost,
		Port:     envs.ZincSearchPort,
		Username: envs.ZincSearchUserName,
		Password: envs.ZincSearchPassword,
	})

	exists, err := zincSearchCli.Index().Exists("user")
	if err != nil {
		if !exists {
			if err = zincSearchCli.Index().Create(&pkgZincSearchModel.Index{
				ShardNum:    1,
				Name:        "user",
				StorageType: "disk",
				Mappings: &pkgZincSearchModel.Mappings{
					Properties: map[string]pkgZincSearchModel.Property{
						"user_id": {
							Type: "text",
						},
						"name_tokens": {
							Type: "keyword",
						},
					},
				},
			}); err != nil {
				pkgLogger.ZapLogger.Logger.Fatal(errors.New("failed to create index: user").Error())
			}
		}
	}

	return &externalSearchEngineClient{
		zincSearch: zincSearchCli,
	}
}
