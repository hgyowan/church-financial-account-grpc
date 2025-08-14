package controller

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	userV1 "github.com/hgyowan/church-financial-account-grpc/gen/user/v1"
	workspaceV1Model "github.com/hgyowan/church-financial-account-grpc/gen/workspace/model/v1"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	"github.com/hgyowan/church-financial-account-grpc/pkg/constant"
	pkgContext "github.com/hgyowan/go-pkg-library/context"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgTime "github.com/hgyowan/go-pkg-library/time"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/samber/lo"
)

func registerUserGRPCHandler(h *grpcHandler) {
	h.UserServiceServer = &userGRPCHandler{h: h}
	userV1.RegisterUserServiceServer(h.externalGRPCServer.Server(), h)
}

type userGRPCHandler struct {
	h *grpcHandler
}

func (u *userGRPCHandler) GetUser(ctx context.Context, request *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	iCtx, err := pkgContext.IncomingContext(ctx).UserID().Scan()
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	res, err := u.h.service.GetUser(ctx, user.GetUserRequest{
		UserID: iCtx.UserID,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.GetUserResponse{
		Id:                res.ID,
		Email:             res.Email,
		Name:              res.Name,
		Nickname:          res.Nickname,
		PhoneNumber:       res.PhoneNumber,
		Provider:          res.Provider,
		IsTermsAgreed:     res.IsTermsAgreed,
		IsMarketingAgreed: res.IsMarketingAgreed,
		Workspaces: lo.Map(res.Workspaces, func(item *workspace.WorkspaceSimple, index int) *workspaceV1Model.WorkspaceSimple {
			return &workspaceV1Model.WorkspaceSimple{
				Id:       item.ID,
				Name:     item.Name,
				IsOwner:  item.IsOwner,
				IsAdmin:  item.IsAdmin,
				JoinedAt: pkgTime.ConvertToSafeTimestamp(item.JoinedAt),
			}
		}),
		RegisteredAt: pkgTime.ConvertToSafeTimestamp(res.RegisteredAt),
	}, nil
}

func (u *userGRPCHandler) RegisterSSOUser(ctx context.Context, request *userV1.RegisterSSOUserRequest) (*userV1.RegisterSSOUserResponse, error) {
	if err := u.h.service.RegisterSSOUser(ctx, user.RegisterSSOUserRequest{
		SocialType:        constant.SocialType(request.GetSocialType()),
		SSOUserID:         request.GetSsoUserId(),
		Name:              request.GetName(),
		Nickname:          request.GetNickname(),
		PhoneNumber:       request.GetPhoneNumber(),
		IsTermsAgreed:     request.GetIsTermsAgreed(),
		IsMarketingAgreed: pkgVariable.ConvertToPointer(request.GetIsMarketingAgreed()),
	}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.RegisterSSOUserResponse{}, nil
}

func (u *userGRPCHandler) LoginSSO(ctx context.Context, request *userV1.LoginSSORequest) (*userV1.LoginSSOResponse, error) {
	iCtx, err := pkgContext.IncomingContext(ctx).IP().UserAgent().Scan()
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	userAgent := internal.ParseUserAgent(iCtx.UserAgent)

	res, err := u.h.service.LoginSSO(ctx, user.LoginSSORequest{
		Code:       request.GetCode(),
		SocialType: constant.SocialType(request.GetSocialType()),
		IP:         iCtx.IP,
		Browser:    userAgent.Browser,
		OS:         userAgent.OS,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.LoginSSOResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (u *userGRPCHandler) LoginEmail(ctx context.Context, request *userV1.LoginEmailRequest) (*userV1.LoginEmailResponse, error) {
	iCtx, err := pkgContext.IncomingContext(ctx).IP().UserAgent().Scan()
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	userAgent := internal.ParseUserAgent(iCtx.UserAgent)

	res, err := u.h.service.LoginEmail(ctx, user.LoginEmailRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		IP:       iCtx.IP,
		Browser:  userAgent.Browser,
		OS:       userAgent.OS,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.LoginEmailResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (u *userGRPCHandler) VerifyEmail(ctx context.Context, request *userV1.VerifyEmailRequest) (*userV1.VerifyEmailResponse, error) {
	if err := u.h.service.VerifyEmail(ctx, user.VerifyEmailRequest{
		Email: request.GetEmail(),
		Code:  request.GetCode(),
	}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.VerifyEmailResponse{}, nil
}

func (u *userGRPCHandler) SendVerifyEmail(ctx context.Context, request *userV1.SendVerifyEmailRequest) (*userV1.SendVerifyEmailResponse, error) {
	if err := u.h.service.SendVerifyEmail(ctx, user.SendVerifyEmailRequest{
		Email: request.GetEmail(),
	}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.SendVerifyEmailResponse{}, nil
}

func (u *userGRPCHandler) RegisterEmailUser(ctx context.Context, request *userV1.RegisterEmailUserRequest) (*userV1.RegisterEmailUserResponse, error) {
	if err := u.h.service.RegisterEmailUser(ctx, user.RegisterEmailUserRequest{
		Name:              request.GetName(),
		Nickname:          request.GetNickname(),
		Email:             request.GetEmail(),
		PhoneNumber:       request.GetPhoneNumber(),
		Password:          request.GetPassword(),
		PasswordConfirm:   request.GetPasswordConfirm(),
		IsTermsAgreed:     request.GetIsTermsAgreed(),
		IsMarketingAgreed: pkgVariable.ConvertToPointer(request.GetIsMarketingAgreed()),
	}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &userV1.RegisterEmailUserResponse{}, nil
}
