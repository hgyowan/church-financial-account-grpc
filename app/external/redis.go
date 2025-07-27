package external

import (
	"context"
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/redis/go-redis/v9"
)

type externalRedisClient struct {
	client *redis.Client
}

func (e *externalRedisClient) Redis() *redis.Client {
	return e.client
}

func MustNewExternalRedis() domain.ExternalRedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", envs.RedisAddr, envs.RedisPort),
		Password: envs.RedisPassword,
		DB:       0,
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		pkgLogger.ZapLogger.Logger.Error(err.Error())
	}
	pkgLogger.ZapLogger.Logger.Info(res)

	return &externalRedisClient{client: client}
}
