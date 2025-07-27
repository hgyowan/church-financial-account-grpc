package service

import "context"

type userService struct {
	s *service
}

func registerUserService(ctx context.Context, s *service) {
	s.UserService = &userService{s: s}
}
