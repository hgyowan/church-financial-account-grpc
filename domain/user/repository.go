package user

type UserRepository interface {
	CreateUser(param *User) error
}
