package controller

import (
	"context"
	userV1 "github.com/hgyowan/church-financial-account-grpc/gen/user/v1"
)

type userGRPCHandler struct {
	h *grpcHandler
}

func registerUserGRPCHandler(ctx context.Context, h *grpcHandler) {
	h.UserServiceServer = &userGRPCHandler{h: h}
	userV1.RegisterUserServiceServer(h.externalGRPCServer.Server(), h)
}
