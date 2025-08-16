package workspace

type WorkspaceRepository interface {
	ListWorkspaceUserByUserID(userID string) ([]*WorkspaceUser, error)
	ListWorkspaceByIDs(ids []string) ([]*Workspace, error)
	PagingWorkspace(param PagingWorkspaceDBParam) ([]*Workspace, bool, error)
	GetWorkspaceByID(id string) (*Workspace, error)
	CreateWorkspaceInvite(param *WorkspaceInvite) error
}
