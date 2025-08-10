package workspace

type WorkspaceRepository interface {
	ListWorkspaceUserByUserID(userID string) ([]*WorkspaceUser, error)
	ListWorkspaceByIDs(ids []string) ([]*Workspace, error)
}
