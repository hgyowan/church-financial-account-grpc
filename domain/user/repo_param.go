package user

import "github.com/hgyowan/church-financial-account-grpc/pkg/constant"

type GetUserSSOByEmailAndProviderAndProviderUserID struct {
	Email      string
	Provider   constant.SocialType
	ProviderID string
}
