package service

import (
	"github.com/google/uuid"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func registerUserService(s *service) {
	s.UserService = &userService{s: s}
}

type userService struct {
	s *service
}

func (u *userService) CreateUser(request user.CreateUserRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	if request.Password != request.PasswordConfirm {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.PasswordMisMatch)
	}

	if request.IsTermsAgreed == false {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.AgreeRequired)
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if request.Nickname == "" {
		request.Nickname = request.Name
	}

	if err = u.s.repo.WithTransaction(func(txRepo domain.Repository) error {
		now := time.Now().UTC()
		userID := uuid.NewString()

		userInfo := &user.User{
			ID:          userID,
			Email:       request.Email,
			Name:        request.Name,
			Nickname:    request.Nickname,
			Password:    string(hashedPw),
			PhoneNumber: request.PhoneNumber,
			CreatedAt:   now,
			UpdatedAt:   &now,
		}

		if err = txRepo.CreateUser(userInfo); err != nil {
			return pkgError.WrapWithCode(err, pkgError.Create)
		}

		if err = txRepo.CreateUserConsent(&user.UserConsent{
			UserID:            userID,
			IsTermsAgreed:     request.IsTermsAgreed,
			IsMarketingAgreed: request.IsMarketingAgreed,
			CreatedAt:         now,
			UpdatedAt:         &now,
		}); err != nil {
			return pkgError.WrapWithCode(err, pkgError.Create)
		}

		return nil
	}); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}
