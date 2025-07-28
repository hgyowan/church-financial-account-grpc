package service

import "github.com/hgyowan/church-financial-account-grpc/domain/user"

func registerUserService(s *service) {
	s.UserService = &userService{s: s}
}

type userService struct {
	s *service
}

func (u userService) CreateUser(request user.CreateUserRequest) error {
	//TODO implement me
	panic("implement me")
}
