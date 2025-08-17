package workspace

import "context"

type WorkspaceService interface {
	ListWorkspaceIntro(ctx context.Context, request ListWorkspaceIntroRequest) (*ListWorkspaceIntroResponse, error)
	SendWorkspaceInviteMessage(ctx context.Context, request SendWorkspaceInviteMessageRequest) error
	ValidWorkspace(ctx context.Context, request ValidWorkspaceRequest) error
}
