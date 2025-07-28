package user

type UserService interface {
	CreateUser(request CreateUserRequest) error
}
