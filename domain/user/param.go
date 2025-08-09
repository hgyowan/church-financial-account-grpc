package user

import "github.com/hgyowan/church-financial-account-grpc/pkg/constant"

type RegisterEmailUserRequest struct {
	Name              string `json:"name" validate:"required"`
	Nickname          string `json:"nickname"`
	Email             string `json:"email" validate:"required,email"`
	PhoneNumber       string `json:"phoneNumber" validate:"phone_number_reg"`
	Password          string `json:"password" validate:"required"`
	PasswordConfirm   string `json:"passwordConfirm" validate:"required"`
	IsTermsAgreed     bool   `json:"isTermsAgreed" validate:"required"`
	IsMarketingAgreed *bool  `json:"isMarketingAgreed" validate:"required"`
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
	SocialType constant.SocialType `json:"socialType" validate:"required"`
	IP         string              `json:"ip"`
	Browser    string              `json:"browser"`
	OS         string              `json:"os"`
}

type LoginSSOResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type IssueTokenRequest struct {
	Code string `json:"code" validate:"required"`
}

type IssueTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GetSSOUserRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

type GetSSOUserResponse struct {
	SSOUser SSOUser `json:"ssoUser"`
}

type RegisterSSOUserRequest struct {
	SocialType        constant.SocialType `json:"socialType" validate:"required"`
	SSOUserID         string              `json:"ssoUserId" validate:"required"`
	Name              string              `json:"name" validate:"required"`
	Nickname          string              `json:"nickname"`
	PhoneNumber       string              `json:"phoneNumber" validate:"omitempty,phone_number_reg"`
	IsTermsAgreed     bool                `json:"isTermsAgreed" validate:"required"`
	IsMarketingAgreed *bool               `json:"isMarketingAgreed" validate:"required"`
}

type RegisterUserRequest struct {
	Email             string              `json:"email"`
	SocialType        constant.SocialType `json:"socialType"`
	SSOUserID         string              `json:"ssoUserId"`
	Password          string              `json:"password"`
	Name              string              `json:"name"`
	Nickname          string              `json:"nickname"`
	PhoneNumber       string              `json:"phoneNumber"`
	IsTermsAgreed     bool                `json:"isTermsAgreed"`
	IsMarketingAgreed *bool               `json:"isMarketingAgreed"`
}

type LoginEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	IP       string `json:"ip"`
	Browser  string `json:"browser"`
	OS       string `json:"os"`
}

type LoginEmailResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
