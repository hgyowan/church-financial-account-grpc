package user

type UserRepository interface {
	CreateUser(param *User) error
	CreateUserConsent(param *UserConsent) error
}
