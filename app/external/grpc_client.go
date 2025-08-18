package external

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	coreWorkspaceGrpc "github.com/hgyowan/church-financial-core-grpc/gen/workspace/v1"
)

type externalGrpcClient struct {
	workspaceClient coreWorkspaceGrpc.WorkspaceServiceClient
}

func (g *externalGrpcClient) WorkspaceClient() coreWorkspaceGrpc.WorkspaceServiceClient {
	return g.workspaceClient
}

func MustNewExternalGrpcClient() domain.ExternalGRPCClient {
	return &externalGrpcClient{
		workspaceClient: coreWorkspaceGrpc.WorkspaceServiceClientProvider(),
	}
}
