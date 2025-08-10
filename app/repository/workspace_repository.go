package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func registerWorkspaceRepository(r *repository) {
	r.WorkspaceRepository = &workspaceRepository{repository: r}
}

type workspaceRepository struct {
	repository *repository
}

func (w *workspaceRepository) ListWorkspaceUserByUserID(userID string) ([]*workspace.WorkspaceUser, error) {
	var res []*workspace.WorkspaceUser
	if err := w.repository.externalGormClient.DB().Where("user_id = ?", userID).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (w *workspaceRepository) ListWorkspaceByIDs(ids []string) ([]*workspace.Workspace, error) {
	var res []*workspace.Workspace
	if err := w.repository.externalGormClient.DB().Where("id IN (?)", ids).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}
