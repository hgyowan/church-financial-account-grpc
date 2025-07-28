package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
)

type repository struct {
	user.UserRepository
	externalGormClient domain.ExternalDBClient
}

func NewRepository(externalGormClient domain.ExternalDBClient) domain.Repository {
	r := &repository{
		externalGormClient: externalGormClient,
	}
	r.register()

	return r
}

func (r *repository) register() {
	registerUserRepository(r)
}

func (r *repository) WithTransaction(fn func(txRepo domain.Repository) error) error {
	tx := r.externalGormClient.DB().Begin() // 트랜잭션 시작
	if tx.Error != nil {
		return tx.Error
	}

	txRepo := &repository{
		externalGormClient: r.externalGormClient.NewTxDB(tx),
	}

	txRepo.register()

	err := fn(txRepo)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
