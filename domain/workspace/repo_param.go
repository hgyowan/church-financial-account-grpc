package workspace

import "github.com/hgyowan/church-financial-account-grpc/internal"

type PagingWorkspaceDBParam struct {
	Sort          internal.WorkspaceSort          `json:"sort"`
	SortDirection internal.WorkspaceSortDirection `json:"sortDirection"`
	PageNo        int                             `json:"pageNo"`
	PageSize      int                             `json:"pageSize"`
}
