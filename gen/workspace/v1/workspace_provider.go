package v1

import (
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLibrary "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
)

func WorkspaceServiceClientProvider() WorkspaceServiceClient {
	conn := pkgLibrary.MustNewGRPCClient(envs.CFMAccountGRPC)
	return NewWorkspaceServiceClient(conn)
}
