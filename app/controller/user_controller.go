package controller

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	userV1 "github.com/hgyowan/church-financial-account-grpc/gen/user/v1"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
)

func registerUserGRPCHandler(h *grpcHandler) {
	h.UserServiceServer = &userGRPCHandler{h: h}
	userV1.RegisterUserServiceServer(h.externalGRPCServer.Server(), h)
}

type userGRPCHandler struct {
	h *grpcHandler
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
