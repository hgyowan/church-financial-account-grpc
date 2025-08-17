package repository

import (
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func registerWorkspaceRepository(r *repository) {
	r.WorkspaceRepository = &workspaceRepository{repository: r}
}

type workspaceRepository struct {
	repository *repository
}

func (w *workspaceRepository) GetWorkspaceByName(name string) (*workspace.Workspace, error) {
	var res *workspace.Workspace
	if err := w.repository.externalGormClient.DB().Where("name = ?", name).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (w *workspaceRepository) CreateWorkspaceInvite(param *workspace.WorkspaceInvite) error {
	return pkgError.Wrap(w.repository.externalGormClient.DB().Create(&param).Error)
}

func (w *workspaceRepository) GetWorkspaceByID(id string) (*workspace.Workspace, error) {
	var res *workspace.Workspace
	if err := w.repository.externalGormClient.DB().Where("id = ?", id).Find(&res).Error; err != nil {
		return nil, pkgError.Wrap(err)
	}
	return res, nil
}

func (w *workspaceRepository) PagingWorkspace(param workspace.PagingWorkspaceDBParam) ([]*workspace.Workspace, bool, error) {
	var res []*workspace.Workspace
	var hasNext bool

	db := w.repository.externalGormClient.DB()

	switch param.Sort {
	case internal.WorkspaceSortName:
		db = db.Order("name " + param.SortDirection)
	case internal.WorkspaceSortMemberCount:
		db = db.Order("member_count " + param.SortDirection)
	}

	if err := db.
		Offset((param.PageNo - 1) * param.PageSize).
		Limit(param.PageSize + 1).
		Find(&res).Error; err != nil {
		return nil, false, pkgError.Wrap(err)
	}

	if len(res) > param.PageSize {
		hasNext = true
		res = res[:param.PageSize]
	}

	return res, hasNext, nil
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
