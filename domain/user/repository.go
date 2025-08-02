package user

type UserRepository interface {
	CreateUser(param *User) error
	CreateUserConsent(param *UserConsent) error
	GetUserByEmail(email string) (*User, error)
	GetUserSSOByEmail(email string) (*UserSSO, error)
	GetUserSSOByEmailAndProviderAndProviderUserID(param GetUserSSOByEmailAndProviderAndProviderUserID) (*UserSSO, error)
}
