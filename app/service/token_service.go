package service

func registerTokenService(s *service) {
	s.TokenService = &tokenService{s: s}
}

type tokenService struct {
	s *service
}
