package external

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	grpcLibrary "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
)

type externalGrpcServer struct {
	server grpcLibrary.GrpcServer
	port   string
}

func (g *externalGrpcServer) Server() grpcLibrary.GrpcServer {
	return g.server
}

func (g *externalGrpcServer) Port() string {
	return g.port
}

func MustNewGRPCServer() domain.ExternalGRPCServer {
	server := grpcLibrary.MustNewGRPCServer()

	return &externalGrpcServer{
		server: server,
		port:   envs.ServerPort,
	}
}
