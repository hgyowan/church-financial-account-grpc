package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/pkg/constant"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_RegisterEmailUser(t *testing.T) {
	beforeEach()
	err := svc.RegisterEmailUser(ctx, user.RegisterEmailUserRequest{
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

func TestUserService_LoginSSO(t *testing.T) {
	beforeEach()
	res, err := svc.LoginSSO(ctx, user.LoginSSORequest{
		Code:       "gWPtJLb0Cl6yYOT2AhbRMAkKzwkoAC74MnMvm_4wXhnw_NrdDLHF5gAAAAQKDSBaAAABmHW5kAf-oZq-Jypvmw",
		SocialType: constant.SocialTypeKakao,
	})
	require.NoError(t, err)
	t.Log(res)
}

func TestUserService_RegisterSSOUser(t *testing.T) {
	beforeEach()
	err := svc.RegisterSSOUser(ctx, user.RegisterSSOUserRequest{
		SocialType:        constant.SocialTypeKakao,
		SSOUserID:         "4378644315",
		Name:              "황교완",
		Nickname:          "임수황",
		PhoneNumber:       "",
		IsTermsAgreed:     true,
		IsMarketingAgreed: pkgVariable.ConvertToPointer(true),
	})
	require.NoError(t, err)
}
