package service

import (
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain/user"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgVariable "github.com/hgyowan/go-pkg-library/variable"
	"github.com/samber/lo"
	"time"
)

func registerWorkspaceService(s *service) {
	s.WorkspaceService = &workspaceService{s: s}
}

type workspaceService struct {
	s *service
}

func (w *workspaceService) SendWorkspaceInviteMessage(request workspace.SendWorkspaceInviteMessageRequest) error {
	if err := w.s.validator.Validator().Struct(request); err != nil {
		return pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	ws, err := w.s.repo.GetWorkspaceByID(request.WorkspaceID)
	if err != nil {
		return pkgError.WrapWithCode(err, pkgError.Get)
	}

	if ws.ID == "" {
		return pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.NotFound)
	}

	if err = w.s.repo.CreateWorkspaceInvite(&workspace.WorkspaceInvite{
		WorkspaceID: request.WorkspaceID,
		UserID:      request.UserID,
		Message:     request.Message,
		CreatedAt:   time.Now().UTC(),
	}); err != nil {
		return pkgError.WrapWithCode(err, pkgError.Create)
	}

	return nil
}

func (w *workspaceService) ListWorkspaceIntro(request workspace.ListWorkspaceIntroRequest) (*workspace.ListWorkspaceIntroResponse, error) {
	if err := w.s.validator.Validator().Struct(request); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	if _, ok := internal.WorkspaceSortMap[request.Sort]; !ok {
		request.Sort = internal.WorkspaceSortName
	}

	if _, ok := internal.WorkspaceSortDirectionMap[request.SortDirection]; !ok {
		request.SortDirection = internal.WorkspaceSortDirectionAsc
	}

	workspaceList, hasNext, err := w.s.repo.PagingWorkspace(workspace.PagingWorkspaceDBParam{
		Sort:          request.Sort,
		SortDirection: request.SortDirection,
		PageNo:        request.PageNo,
		PageSize:      request.PageSize,
	})
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	ownerIDs := lo.Map(workspaceList, func(item *workspace.Workspace, index int) string {
		return item.OwnerID
	})

	userList, err := w.s.repo.ListUserByIDs(ownerIDs)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	userMap := lo.SliceToMap(userList, func(item *user.User) (string, *user.User) {
		return item.ID, item
	})

	return &workspace.ListWorkspaceIntroResponse{
		List: lo.Map(workspaceList, func(item *workspace.Workspace, index int) *workspace.WorkspaceIntro {
			ownerName := ""
			if u, ok := userMap[item.OwnerID]; ok {
				ownerName = u.Name
			}
			return &workspace.WorkspaceIntro{
				ThumbnailURL: pkgVariable.GetSafeValue(item.ThumbnailURL, ""),
				Name:         item.Name,
				Address:      fmt.Sprintf("%s %s", item.Address1, item.Address2),
				OwnerName:    ownerName,
				Description:  pkgVariable.GetSafeValue(item.Description, ""),
				MemberCount:  item.MemberCount,
			}
		}),
		HasNext: hasNext,
	}, nil
}
