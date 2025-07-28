package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	beforeEach()
	err := svc.CreateUser(user.CreateUserRequest{
		Name:                  "황교완",
		Nickname:              "임수황",
		Email:                 "test@gmail.com",
		EmailVerifyCode:       "10001",
		PhoneNumber:           "010-1234-1234",
		PhoneNumberVerifyCode: "50681",
		Password:              "test",
		PasswordConfirm:       "test",
		IsTermsAgreed:         true,
		IsMarketingAgreed:     pkgVariable.ConvertToPointer(true),
	})
	require.NoError(t, err)
}
