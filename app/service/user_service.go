package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	"github.com/hgyowan/church-financial-account-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
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

func (u *userService) RegisterSSOUser(ctx context.Context, request user.RegisterSSOUserRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	tk, err := u.s.externalRedisClient.Redis().Get(ctx, fmt.Sprintf("sso:%s:id:%s", request.SocialType, request.SSOUserID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return pkgError.WrapWithCode(err, pkgError.Expired)
		}

		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	sso, err := NewSSOService(u.s, request.SocialType)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	ssoUser, err := sso.GetSSOUser(ctx, user.GetSSOUserRequest{AccessToken: tk})
	if err != nil {
		return pkgError.Wrap(err)
	}

	// 중복 체크
	userData, err := u.s.repo.GetUserByEmail(ssoUser.SSOUser.Email)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(ssoUser.SSOUser.Email)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	if request.IsTermsAgreed == false {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.AgreeRequired)
	}

	if request.Nickname == "" && ssoUser.SSOUser.Nickname != "" {
		request.Nickname = ssoUser.SSOUser.Nickname
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
			Email:       ssoUser.SSOUser.Email,
			Name:        request.Name,
			Nickname:    request.Nickname,
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

		if err = txRepo.CreateUserSSO(&user.UserSSO{
			UserID:         userID,
			Provider:       string(request.SocialType),
			ProviderUserID: ssoUser.SSOUser.SSOUserID,
			Email:          ssoUser.SSOUser.Email,
			CreatedAt:      now,
		}); err != nil {
			return pkgError.WrapWithCode(err, pkgError.Create)
		}

		return nil
	}); err != nil {
		return pkgError.Wrap(err)
	}

	if err := u.s.externalRedisClient.Redis().Del(ctx, fmt.Sprintf("sso:%s:id:%s", request.SocialType, request.SSOUserID)).Err(); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	// TODO: zincsearch 형태소 (hmac) 저장

	return nil
}

func (u *userService) LoginSSO(ctx context.Context, request user.LoginSSORequest) (*user.LoginSSOResponse, error) {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	sso, err := NewSSOService(u.s, request.SocialType)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	tk, err := sso.IssueToken(ctx, user.IssueTokenRequest{Code: request.Code})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	ssoUser, err := sso.GetSSOUser(ctx, user.GetSSOUserRequest{AccessToken: tk.AccessToken})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	userSSO, err := u.s.repo.GetUserSSOByEmailAndProviderAndProviderUserID(user.GetUserSSOByEmailAndProviderAndProviderUserID{
		Email:      ssoUser.SSOUser.Email,
		Provider:   request.SocialType,
		ProviderID: ssoUser.SSOUser.SSOUserID,
	})
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userSSO.UserID == "" {
		// sso_user_id 와 access_token 을 redis 에 30분 동안 캐싱하고 회원가입에 활용한다.
		if err = u.s.externalRedisClient.Redis().Set(ctx, fmt.Sprintf("sso:%s:id:%s", request.SocialType, ssoUser.SSOUser.SSOUserID), tk.AccessToken, time.Minute*30).Err(); err != nil {
			return nil, pkgError.WrapWithCode(err, pkgError.Create)
		}

		return nil, pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.NotFound, ssoUser.SSOUser.SSOUserID)
	}

	jwtToken, err := u.s.IssueJWTToken(ctx, token.IssueJWTTokenRequest{
		UserID: userSSO.UserID,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if err = u.s.repo.CreateUserLoginLog(&user.UserLoginLog{
		UserID:    userSSO.UserID,
		IP:        pkgVariable.ConvertToPointer(request.IP),
		Browser:   pkgVariable.ConvertToPointer(request.Browser),
		OS:        pkgVariable.ConvertToPointer(request.OS),
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	return &user.LoginSSOResponse{
		AccessToken:  jwtToken.JWTToken.AccessToken,
		RefreshToken: jwtToken.JWTToken.RefreshToken,
	}, nil
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

func (u *userService) RegisterEmailUser(ctx context.Context, request user.RegisterEmailUserRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	// 중복 체크
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

		if err = txRepo.CreateUserSSO(&user.UserSSO{
			UserID:         userID,
			Provider:       string(constant.SocialTypeEmail),
			ProviderUserID: "",
			Email:          request.Email,
			CreatedAt:      now,
		}); err != nil {
			return pkgError.WrapWithCode(err, pkgError.Create)
		}

		return nil
	}); err != nil {
		return pkgError.Wrap(err)
	}

	// TODO: zincsearch 형태소 (hmac) 저장

	return nil
}
