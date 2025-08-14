package internal

type WorkspaceSort string

const (
	WorkspaceSortName        WorkspaceSort = "name"
	WorkspaceSortMemberCount WorkspaceSort = "memberCount"
)

var WorkspaceSortMap = map[WorkspaceSort]struct{}{
	WorkspaceSortName:        {},
	WorkspaceSortMemberCount: {},
}

type WorkspaceSortDirection string

const (
	WorkspaceSortDirectionAsc  WorkspaceSortDirection = "asc"
	WorkspaceSortDirectionDesc WorkspaceSortDirection = "desc"
)

var WorkspaceSortDirectionMap = map[WorkspaceSortDirection]struct{}{
	WorkspaceSortDirectionAsc:  {},
	WorkspaceSortDirectionDesc: {},
}
