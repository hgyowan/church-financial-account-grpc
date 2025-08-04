package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func registerUserRepository(r *repository) {
	r.UserRepository = &userRepository{repository: r}
}

type userRepository struct {
	repository *repository
}

func (u *userRepository) CreateUserLoginLog(param *user.UserLoginLog) error {
	return pkgError.Wrap(u.repository.externalGormClient.DB().Create(&param).Error)
}

func (u *userRepository) CreateUserSSO(param *user.UserSSO) error {
	return pkgError.Wrap(u.repository.externalGormClient.DB().Create(&param).Error)
}

func (u *userRepository) GetUserSSOByEmailAndProviderAndProviderUserID(param user.GetUserSSOByEmailAndProviderAndProviderUserID) (*user.UserSSO, error) {
	var res *user.UserSSO
	if err := u.repository.externalGormClient.DB().
		Where("email = ?", param.Email).
		Where("provider = ?", param.Provider).
		Where("provider_user_id = ?", param.ProviderID).
		Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (u *userRepository) GetUserSSOByEmail(email string) (*user.UserSSO, error) {
	var res *user.UserSSO
	if err := u.repository.externalGormClient.DB().Where("email = ?", email).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (u *userRepository) GetUserByEmail(email string) (*user.User, error) {
	var res *user.User
	if err := u.repository.externalGormClient.DB().Where("email = ?", email).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (u *userRepository) CreateUserConsent(param *user.UserConsent) error {
	return pkgError.Wrap(u.repository.externalGormClient.DB().Create(&param).Error)
}

func (u *userRepository) CreateUser(param *user.User) error {
	return pkgError.Wrap(u.repository.externalGormClient.DB().Create(&param).Error)
}
