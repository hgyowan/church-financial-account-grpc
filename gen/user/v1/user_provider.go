package v1

import (
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLibrary "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
)

func UserServiceClientProvider() UserServiceClient {
	conn := pkgLibrary.MustNewGRPCClient(envs.CFMAccountGRPC)
	return NewUserServiceClient(conn)
}
