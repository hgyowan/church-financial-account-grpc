package user

type CreateEmailUserRequest struct {
	Name              string `json:"name" validate:"required"`
	Nickname          string `json:"nickname"`
	Email             string `json:"email" validate:"required,email"`
	PhoneNumber       string `json:"phone_number" validate:"required,phone_number_reg"`
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
