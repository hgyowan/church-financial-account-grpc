package workspace

import "github.com/hgyowan/church-financial-account-grpc/internal"

type ListWorkspaceIntroRequest struct {
	Sort          internal.WorkspaceSort          `json:"sort"`
	SortDirection internal.WorkspaceSortDirection `json:"sortDirection"`
	PageNo        int                             `json:"pageNo" validate:"required"`
	PageSize      int                             `json:"pageSize" validate:"required"`
}

type ListWorkspaceIntroResponse struct {
	List    []*WorkspaceIntro `json:"list"`
	HasNext bool              `json:"hasNext"`
}

type SendWorkspaceInviteMessageRequest struct {
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Message     string `json:"message" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
}
