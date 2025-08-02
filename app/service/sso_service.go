package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func NewSSOService(s *service, socialType constant.SocialType) (user.SSOService, error) {
	switch socialType {
	case constant.SocialTypeKakao:
		return NewSSOServiceKakao(s), nil
	}

	return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.UnsupportedOAuthProvider)
}
