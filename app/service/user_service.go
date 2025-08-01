package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func registerUserService(s *service) {
	s.UserService = &userService{s: s}
}

type userService struct {
	s *service
}

func (u *userService) VerifyEmail(ctx context.Context, request user.VerifyEmailRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	userData, err := u.s.repo.GetUserByEmail(request.Email)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(request.Email)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	code, err := u.s.externalRedisClient.Redis().Get(ctx, fmt.Sprintf("emailVerify:%s", request.Email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return pkgError.WrapWithCode(err, pkgError.Expired)
		}

		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if request.Code != code {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.WrongParam)
	}

	return nil
}

func (u *userService) SendVerifyEmail(ctx context.Context, request user.SendVerifyEmailRequest) error {
	userData, err := u.s.repo.GetUserByEmail(request.Email)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(request.Email)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	code := internal.GenerateRandomCode()
	if err = u.s.externalRedisClient.Redis().Set(ctx, fmt.Sprintf("emailVerify:%s", request.Email), code, time.Minute*30).Err(); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Create)
	}

	type templateData struct {
		VerifyCode string
	}

	data := templateData{
		VerifyCode: code,
	}

	if err = u.s.externalMailSender.MailSender().SendMailWithTemplate(request.Email, "[CFM] 회원가입 이메일 인증", pkgEmail.TemplateKeyVerifyEmail, data); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (u *userService) CreateEmailUser(ctx context.Context, request user.CreateEmailUserRequest) error {
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

	// TODO: 암호화
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
