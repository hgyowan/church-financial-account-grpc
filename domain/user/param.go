package user

type CreateUserRequest struct {
	Name            string `json:"name"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	EmailVerifyCode string `json:"email_verify_code"`
}
