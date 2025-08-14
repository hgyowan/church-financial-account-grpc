package workspace

type WorkspaceService interface {
	ListWorkspaceIntro(request ListWorkspaceIntroRequest) (*ListWorkspaceIntroResponse, error)
	SendWorkspaceInviteMessage(request SendWorkspaceInviteMessageRequest) error
}
