package service

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_ListWorkspaceIntro(t *testing.T) {
	beforeEach()
	res, err := svc.ListWorkspaceIntro(ctx, workspace.ListWorkspaceIntroRequest{
		Sort:          "",
		SortDirection: "",
		PageNo:        1,
		PageSize:      10,
	})
	require.NoError(t, err)
	t.Log(res)
}

func TestService_SendWorkspaceInviteMessage(t *testing.T) {
	beforeEach()
	err := svc.SendWorkspaceInviteMessage(ctx, workspace.SendWorkspaceInviteMessageRequest{
		WorkspaceID: "0f5b0151-b4c8-44d1-981a-a25bd426f1e2",
		UserID:      "0f5b0151-b4c8-44d1-981a-a25bd426f1ec",
		Message:     "소개입니다",
	})
	require.NoError(t, err)
}

func TestService_ValidWorkspace(t *testing.T) {
	beforeEach()
	err := svc.ValidWorkspace(ctx, workspace.ValidWorkspaceRequest{Name: "심심교회"})
	require.NoError(t, err)
}
