package external

import (
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	maxOpenConnNum  = 25
	maxIdleConnNum  = 25
	connMaxLifetime = 5 * time.Minute
)

type externalDB struct {
	client *gorm.DB
}

func (e *externalDB) DB() *gorm.DB {
	return e.client
}

func (e *externalDB) NewTxDB(tx *gorm.DB) domain.ExternalDBClient {
	return &externalDB{client: tx}
}

func MustNewExternalDB() domain.ExternalDBClient {
	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		envs.DBHost, envs.DBUser, envs.DBPassword, envs.DBName, envs.DBPort,
	)

	// Initialize GORM DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 pkgLogger.ZapLogger.GormLogger,
	})
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("failed to connect to database: %v", err)
	}

	// Configure a connection pool
	sqlDB, err := db.DB()
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("failed to get database object: %v", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConnNum)
	sqlDB.SetMaxIdleConns(maxIdleConnNum)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	pkgLogger.ZapLogger.Logger.Info("PostgreSQL connection established successfully!")

	return &externalDB{client: db}
}
