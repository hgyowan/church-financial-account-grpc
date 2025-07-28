package controller

import (
	userV1 "github.com/hgyowan/church-financial-account-grpc/gen/user/v1"
)

func registerUserGRPCHandler(h *grpcHandler) {
	h.UserServiceServer = &userGRPCHandler{h: h}
	userV1.RegisterUserServiceServer(h.externalGRPCServer.Server(), h)
}

type userGRPCHandler struct {
	h *grpcHandler
}
