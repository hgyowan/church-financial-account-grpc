package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_CreateEmailUser(t *testing.T) {
	beforeEach()
	err := svc.CreateEmailUser(ctx, user.CreateEmailUserRequest{
		Name:              "황교완",
		Nickname:          "임수황",
		Email:             "test@gmail.com",
		PhoneNumber:       "010-1234-1234",
		Password:          "test",
		PasswordConfirm:   "test",
		IsTermsAgreed:     true,
		IsMarketingAgreed: pkgVariable.ConvertToPointer(true),
	})
	require.NoError(t, err)
}

func TestUserService_SendVerifyEmail(t *testing.T) {
	beforeEach()
	err := svc.SendVerifyEmail(ctx, user.SendVerifyEmailRequest{
		Email: "rydhkstptkd@naver.com",
	})
	require.NoError(t, err)
}

func TestUserService_VerifyEmail(t *testing.T) {
	beforeEach()
	err := svc.VerifyEmail(ctx, user.VerifyEmailRequest{
		Email: "rydhkstptkd@naver.com",
		Code:  "893359",
	})
	require.NoError(t, err)
}
