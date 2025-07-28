package main

import (
	"context"
	grpc "github.com/hgyowan/church-financial-account-grpc/app/controller"
	"github.com/hgyowan/church-financial-account-grpc/app/external"
	"github.com/hgyowan/church-financial-account-grpc/app/repository"
	"github.com/hgyowan/church-financial-account-grpc/app/service"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pkgLogger.MustInitZapLogger()

	if pkgLogger.ZapLogger == nil {
		log.Fatal("logger is nil")
	}

	bCtx, cancelFunc := context.WithCancel(context.Background())
	group, gCtx := errgroup.WithContext(bCtx)
	doneChan := make(chan struct{}, 1)
	grpcServer := external.MustNewGRPCServer()
	dbClient := external.MustNewExternalDB()
	repo := repository.NewRepository(dbClient)
	redisCli := external.MustNewExternalRedis()
	v := external.MustNewValidator()
	svc := service.NewService(repo, redisCli, v)
	pkgLogger.ZapLogger.Logger.Info("Starting gRPC server on")
	group.Go(func() error {
		grpc.NewGRPCHandler(svc, grpcServer).Listen(gCtx)
		pkgLogger.ZapLogger.Logger.Fatal("GRPC Handler End")
		doneChan <- struct{}{}
		return nil
	})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer close(interrupt)

	select {
	case <-doneChan:
		cancelFunc()
	case <-interrupt:
		cancelFunc()
	}

	if err := group.Wait(); err != nil {
		pkgLogger.ZapLogger.Logger.Fatal(err.Error())
	}

	pkgLogger.ZapLogger.Logger.Info("GRPC Server End")
}
