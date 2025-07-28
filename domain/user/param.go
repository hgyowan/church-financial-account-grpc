package user

type CreateUserRequest struct {
	Name                  string `json:"name" validate:"required"`
	Nickname              string `json:"nickname"`
	Email                 string `json:"email" validate:"required,email"`
	EmailVerifyCode       string `json:"email_verify_code" validate:"required"`
	PhoneNumber           string `json:"phone_number" validate:"required,phone_number_reg"`
	PhoneNumberVerifyCode string `json:"phone_number_verify_code" validate:"required"`
	Password              string `json:"password" validate:"required"`
	PasswordConfirm       string `json:"password_confirm" validate:"required"`
	IsTermsAgreed         bool   `json:"is_terms_agreed" validate:"required"`
	IsMarketingAgreed     *bool  `json:"is_marketing_agreed" validate:"required"`
}
