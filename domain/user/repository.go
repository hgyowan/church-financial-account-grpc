package user

type UserRepository interface {
	CreateUser(param *User) error
	CreateUserConsent(param *UserConsent) error
	GetUserByEmail(email string) (*User, error)
	GetUserSSOByEmail(email string) (*UserSSO, error)
	GetUserSSOByEmailAndProviderAndProviderUserID(param GetUserSSOByEmailAndProviderAndProviderUserID) (*UserSSO, error)
	CreateUserSSO(param *UserSSO) error
	CreateUserLoginLog(param *UserLoginLog) error
	GetUserByID(id string) (*User, error)
	GetUserConsentByID(id string) (*UserConsent, error)
	GetUserSSOByID(id string) (*UserSSO, error)
}
