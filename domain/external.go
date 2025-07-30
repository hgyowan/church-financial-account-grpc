package domain

import (
	"github.com/go-playground/validator/v10"
	pkgGrpc "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ExternalGRPCServer interface {
	Server() pkgGrpc.GrpcServer
	Port() string
}

type ExternalDBClient interface {
	DB() *gorm.DB
	NewTxDB(tx *gorm.DB) ExternalDBClient
}

type ExternalRedisClient interface {
	Redis() *redis.Client
}

type ExternalValidator interface {
	Validator() *validator.Validate
}

type ExternalMailSender interface {
	MailSender() pkgEmail.EmailSender
}
