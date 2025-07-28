package controller

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	userV1 "github.com/hgyowan/church-financial-account-grpc/gen/user/v1"
)

type grpcHandler struct {
	userV1.UserServiceServer
	service            domain.Service
	externalGRPCServer domain.ExternalGRPCServer
}

func NewGRPCHandler(service domain.Service, externalGRPCServer domain.ExternalGRPCServer) *grpcHandler {
	h := &grpcHandler{
		service:            service,
		externalGRPCServer: externalGRPCServer,
	}

	h.register()

	return h
}

func (h *grpcHandler) Listen(ctx context.Context) {
	h.externalGRPCServer.Server().Serve(ctx, h.externalGRPCServer.Port())
}

func (h *grpcHandler) register() {
	registerUserGRPCHandler(h)
}
