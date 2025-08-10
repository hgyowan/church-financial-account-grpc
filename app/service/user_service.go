package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	"github.com/hgyowan/church-financial-account-grpc/pkg/constant"
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	pkgNgram "github.com/hgyowan/go-pkg-library/ngram"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	pkgZincSearchModel "github.com/hgyowan/go-pkg-library/zincsearch/model"
	pkgZincSearchParam "github.com/hgyowan/go-pkg-library/zincsearch/param"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func registerUserService(s *service) {
	s.UserService = &userService{s: s}
}

type userService struct {
	s *service
}

func (u *userService) GetUser(ctx context.Context, request user.GetUserRequest) (*user.GetUserResponse, error) {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	userInfo, err := u.s.repo.GetUserByID(request.UserID)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userInfo.ID == "" {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.NotFound)
	}

	userSSO, err := u.s.repo.GetUserSSOByEmail(request.UserID)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	userConsent, err := u.s.repo.GetUserConsentByID(request.UserID)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	workspaceUserList, err := u.s.repo.ListWorkspaceUserByUserID(request.UserID)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	workspaceIDs := lo.Map(workspaceUserList, func(item *workspace.WorkspaceUser, index int) string {
		return item.WorkspaceID
	})

	workspaceList, err := u.s.repo.ListWorkspaceByIDs(workspaceIDs)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	workspaceMap := lo.SliceToMap(workspaceList, func(item *workspace.Workspace) (string, *workspace.Workspace) {
		return item.ID, item
	})

	return &user.GetUserResponse{
		ID:                userInfo.ID,
		Email:             userInfo.Email,
		Name:              userInfo.Name,
		Nickname:          userInfo.Nickname,
		PhoneNumber:       userInfo.PhoneNumber,
		Provider:          userSSO.Provider,
		IsTermsAgreed:     userConsent.IsTermsAgreed,
		IsMarketingAgreed: pkgVariable.GetSafeValue(userConsent.IsMarketingAgreed, false),
		Workspaces: lo.Map(workspaceUserList, func(item *workspace.WorkspaceUser, index int) *user.Workspace {
			if w, ok := workspaceMap[item.WorkspaceID]; ok {
				return &user.Workspace{
					ID:       w.ID,
					Name:     w.Name,
					IsOwner:  w.OwnerID == request.UserID,
					IsAdmin:  item.IsAdmin,
					JoinedAt: item.CreatedAt,
				}
			}

			return nil
		}),
		RegisteredAt: userInfo.CreatedAt,
	}, nil
}

func (u *userService) LoginEmail(ctx context.Context, request user.LoginEmailRequest) (*user.LoginEmailResponse, error) {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	encryptEmail, err := pkgCrypto.CBCEncryptWithFixedKey(request.Email, "email", []byte(envs.MasterKey))
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	userInfo, err := u.s.repo.GetUserByEmail(encryptEmail)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userInfo.ID == "" {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.NotFound)
	}

	userSSOData, err := u.s.repo.GetUserSSOByEmail(encryptEmail)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userSSOData.UserID != "" && userSSOData.Provider != string(constant.SocialTypeEmail) {
		return nil, pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.AlreadyExistsEmail, userSSOData.Provider)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(request.Password)); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.InvalidPassword)
	}

	jwtToken, err := u.s.IssueJWTToken(ctx, token.IssueJWTTokenRequest{
		UserID: userInfo.ID,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	if err = u.s.repo.CreateUserLoginLog(&user.UserLoginLog{
		UserID:    userInfo.ID,
		IP:        pkgVariable.ConvertToPointer(request.IP),
		Browser:   pkgVariable.ConvertToPointer(request.Browser),
		OS:        pkgVariable.ConvertToPointer(request.OS),
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	return &user.LoginEmailResponse{
		AccessToken:  jwtToken.JWTToken.AccessToken,
		RefreshToken: jwtToken.JWTToken.RefreshToken,
	}, nil
}

func (u *userService) registerUser(ctx context.Context, txRepo domain.Repository, request user.RegisterUserRequest) error {
	// 중복 체크
	encryptEmail, err := pkgCrypto.CBCEncryptWithFixedKey(request.Email, "email", []byte(envs.MasterKey))
	if err != nil {
		return pkgError.Wrap(err)
	}

	userData, err := u.s.repo.GetUserByEmail(encryptEmail)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(encryptEmail)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	if request.IsTermsAgreed == false {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.AgreeRequired)
	}

	userID := uuid.NewString()
	now := time.Now().UTC()
	userInfo := &user.User{
		ID:          userID,
		Email:       request.Email,
		Name:        request.Name,
		Nickname:    request.Nickname,
		Password:    request.Password,
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
		ProviderUserID: request.SSOUserID,
		Email:          request.Email,
		CreatedAt:      now,
	}); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Create)
	}

	tokens := pkgNgram.GenerateHmacTokens(request.Name)
	if err = u.s.externalSearchEngine.ZincSearch().Document("user").Create(&pkgZincSearchParam.DocumentCreateRequest{
		ID: userID,
		Document: &pkgZincSearchModel.Document{
			Key: "name",
			Val: tokens,
		},
	}); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Create)
	}

	return nil
}

func (u *userService) RegisterSSOUser(ctx context.Context, request user.RegisterSSOUserRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	tk, err := u.s.externalRedisClient.Redis().Get(ctx, fmt.Sprintf("%s:sso:%s:id:%s", envs.ServiceType, request.SocialType, request.SSOUserID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return pkgError.WrapWithCode(err, pkgError.Expired)
		}

		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	sso, err := NewSSOService(u.s, request.SocialType)
	if err != nil {
		return pkgError.Wrap(err)
	}

	ssoUser, err := sso.GetSSOUser(ctx, user.GetSSOUserRequest{AccessToken: tk})
	if err != nil {
		return pkgError.Wrap(err)
	}

	if request.Nickname == "" && ssoUser.SSOUser.Nickname != "" {
		request.Nickname = ssoUser.SSOUser.Nickname
	}

	if request.Nickname == "" {
		request.Nickname = request.Name
	}

	if err = u.s.repo.WithTransaction(func(txRepo domain.Repository) error {
		if err = u.registerUser(ctx, txRepo, user.RegisterUserRequest{
			Email:             ssoUser.SSOUser.Email,
			Name:              request.Name,
			Nickname:          request.Nickname,
			Password:          "",
			PhoneNumber:       request.PhoneNumber,
			IsTermsAgreed:     request.IsTermsAgreed,
			IsMarketingAgreed: request.IsMarketingAgreed,
			SSOUserID:         request.SSOUserID,
			SocialType:        request.SocialType,
		}); err != nil {
			return pkgError.Wrap(err)
		}

		return nil
	}); err != nil {
		return pkgError.Wrap(err)
	}

	if err := u.s.externalRedisClient.Redis().Del(ctx, fmt.Sprintf("%s:sso:%s:id:%s", envs.ServiceType, request.SocialType, request.SSOUserID)).Err(); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	return nil
}

func (u *userService) LoginSSO(ctx context.Context, request user.LoginSSORequest) (*user.LoginSSOResponse, error) {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	sso, err := NewSSOService(u.s, request.SocialType)
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	tk, err := sso.IssueToken(ctx, user.IssueTokenRequest{Code: request.Code})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	ssoUser, err := sso.GetSSOUser(ctx, user.GetSSOUserRequest{AccessToken: tk.AccessToken})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	encryptEmail, err := pkgCrypto.CBCEncryptWithFixedKey(ssoUser.SSOUser.Email, "email", []byte(envs.MasterKey))
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	userSSOData, err := u.s.repo.GetUserSSOByEmail(encryptEmail)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userSSOData.UserID != "" && userSSOData.Provider != string(request.SocialType) {
		return nil, pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.AlreadyExistsEmail, userSSOData.Provider)
	}

	userSSO, err := u.s.repo.GetUserSSOByEmailAndProviderAndProviderUserID(user.GetUserSSOByEmailAndProviderAndProviderUserID{
		Email:      encryptEmail,
		Provider:   request.SocialType,
		ProviderID: ssoUser.SSOUser.SSOUserID,
	})
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userSSO.UserID == "" {
		// sso_user_id 와 access_token 을 redis 에 30분 동안 캐싱하고 회원가입에 활용한다.
		if err = u.s.externalRedisClient.Redis().Set(ctx, fmt.Sprintf("%s:sso:%s:id:%s", envs.ServiceType, request.SocialType, ssoUser.SSOUser.SSOUserID), tk.AccessToken, time.Minute*30).Err(); err != nil {
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

	encryptEmail, err := pkgCrypto.CBCEncryptWithFixedKey(request.Email, "email", []byte(envs.MasterKey))
	if err != nil {
		return pkgError.Wrap(err)
	}

	userData, err := u.s.repo.GetUserByEmail(encryptEmail)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(encryptEmail)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	code, err := u.s.externalRedisClient.Redis().Get(ctx, fmt.Sprintf("%s:emailVerify:%s", envs.ServiceType, request.Email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return pkgError.WrapWithCode(err, pkgError.Expired)
		}

		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if request.Code != code {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.WrongParam)
	}

	if err = u.s.externalRedisClient.Redis().Del(ctx, fmt.Sprintf("%s:emailVerify:%s", envs.ServiceType, request.Email)).Err(); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Delete)
	}

	return nil
}

func (u *userService) SendVerifyEmail(ctx context.Context, request user.SendVerifyEmailRequest) error {
	if err := u.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	encryptEmail, err := pkgCrypto.CBCEncryptWithFixedKey(request.Email, "email", []byte(envs.MasterKey))
	if err != nil {
		return pkgError.Wrap(err)
	}

	userData, err := u.s.repo.GetUserByEmail(encryptEmail)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if userData.ID != "" {
		userSSOData, err := u.s.repo.GetUserSSOByEmail(encryptEmail)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		return pkgError.WrapWithCodeAndData(pkgError.EmptyBusinessError(), pkgError.Duplicate, userSSOData.Provider)
	}

	code := internal.GenerateRandomCode()
	if err = u.s.externalRedisClient.Redis().Set(ctx, fmt.Sprintf("%s:emailVerify:%s", envs.ServiceType, request.Email), code, time.Minute*30).Err(); err != nil {
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

	if request.Password != request.PasswordConfirm {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.PasswordMisMatch)
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if request.Nickname == "" {
		request.Nickname = request.Name
	}

	if err = u.s.repo.WithTransaction(func(txRepo domain.Repository) error {
		if err = u.registerUser(ctx, txRepo, user.RegisterUserRequest{
			Email:             request.Email,
			Name:              request.Name,
			Nickname:          request.Nickname,
			Password:          string(hashedPw),
			PhoneNumber:       request.PhoneNumber,
			IsTermsAgreed:     request.IsTermsAgreed,
			IsMarketingAgreed: request.IsMarketingAgreed,
			SSOUserID:         "",
			SocialType:        constant.SocialTypeEmail,
		}); err != nil {
			return pkgError.Wrap(err)
		}

		return nil
	}); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}
