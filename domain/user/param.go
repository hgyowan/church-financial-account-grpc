package user

import "github.com/hgyowan/church-financial-account-grpc/pkg/constant"

type RegisterEmailUserRequest struct {
	Name              string `json:"name" validate:"required"`
	Nickname          string `json:"nickname"`
	Email             string `json:"email" validate:"required,email"`
	PhoneNumber       string `json:"phone_number" validate:"phone_number_reg"`
	Password          string `json:"password" validate:"required"`
	PasswordConfirm   string `json:"password_confirm" validate:"required"`
	IsTermsAgreed     bool   `json:"is_terms_agreed" validate:"required"`
	IsMarketingAgreed *bool  `json:"is_marketing_agreed" validate:"required"`
}

type SendVerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type LoginSSORequest struct {
	Code       string              `json:"code" validate:"required"`
	SocialType constant.SocialType `json:"social_type" validate:"required"`
}

type LoginSSOResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type IssueTokenRequest struct {
	Code string `json:"code" validate:"required"`
}

type IssueTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetSSOUserRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type GetSSOUserResponse struct {
	SSOUser SSOUser `json:"sso_user"`
}

type RegisterSSOUserRequest struct {
	SocialType        constant.SocialType `json:"social_type" validate:"required"`
	SSOUserID         string              `json:"sso_user_id" validate:"required"`
	Name              string              `json:"name" validate:"required"`
	Nickname          string              `json:"nickname"`
	PhoneNumber       string              `json:"phone_number" validate:"omitempty,phone_number_reg"`
	IsTermsAgreed     bool                `json:"is_terms_agreed" validate:"required"`
	IsMarketingAgreed *bool               `json:"is_marketing_agreed" validate:"required"`
}
