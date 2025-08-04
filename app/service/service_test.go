package service

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/app/client"
	"github.com/hgyowan/church-financial-account-grpc/app/external"
	"github.com/hgyowan/church-financial-account-grpc/app/repository"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
)

var ctx context.Context
var svc domain.Service

func beforeEach() {
	pkgLogger.MustInitZapLogger()
	db := external.MustNewExternalDB()
	redis := external.MustNewExternalRedis()
	v := external.MustNewValidator()
	mailSender := external.MustNewEmailSender("/Users/hwang-gyowan/go/src/church-financial-account-grpc/internal/format/")
	repo := repository.NewRepository(db)
	http := external.MustNewExternalHttpClient()
	cli := client.NewClient(http)
	svc = NewService(repo, cli, redis, mailSender, v)
	ctx = context.Background()
}
