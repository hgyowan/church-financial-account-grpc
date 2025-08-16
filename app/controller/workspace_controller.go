package controller

import (
	"context"
	"github.com/hgyowan/church-financial-account-grpc/domain/workspace"
	workspaceV1Model "github.com/hgyowan/church-financial-account-grpc/gen/workspace/model/v1"
	workspaceV1 "github.com/hgyowan/church-financial-account-grpc/gen/workspace/v1"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	pkgContext "github.com/hgyowan/go-pkg-library/context"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"github.com/samber/lo"
)

func registerWorkspaceGRPCHandler(h *grpcHandler) {
	h.WorkspaceServiceServer = &workspaceGRPCHandler{h: h}
	workspaceV1.RegisterWorkspaceServiceServer(h.externalGRPCServer.Server(), h)
}

type workspaceGRPCHandler struct {
	h *grpcHandler
}

func (w *workspaceGRPCHandler) ListWorkspaceIntro(ctx context.Context, request *workspaceV1.ListWorkspaceIntroRequest) (*workspaceV1.ListWorkspaceIntroResponse, error) {
	res, err := w.h.service.ListWorkspaceIntro(ctx, workspace.ListWorkspaceIntroRequest{
		Sort:          internal.WorkspaceSort(request.GetSort()),
		SortDirection: internal.WorkspaceSortDirection(request.GetSortDirection()),
		PageNo:        int(request.GetPageNo()),
		PageSize:      int(request.GetPageSize()),
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &workspaceV1.ListWorkspaceIntroResponse{
		List: lo.Map(res.List, func(item *workspace.WorkspaceIntro, index int) *workspaceV1Model.WorkspaceIntro {
			return &workspaceV1Model.WorkspaceIntro{
				Id:           item.ID,
				ThumbnailUrl: item.ThumbnailURL,
				Name:         item.Name,
				Address:      item.Address,
				OwnerName:    item.OwnerName,
				Description:  item.Description,
				MemberCount:  int32(item.MemberCount),
			}
		}),
		HasNext: res.HasNext,
	}, nil
}

func (w *workspaceGRPCHandler) SendWorkspaceInviteMessage(ctx context.Context, request *workspaceV1.SendWorkspaceInviteMessageRequest) (*workspaceV1.SendWorkspaceInviteMessageResponse, error) {
	iCtx, err := pkgContext.IncomingContext(ctx).UserID().Scan()
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
	}

	if err = w.h.service.SendWorkspaceInviteMessage(ctx, workspace.SendWorkspaceInviteMessageRequest{
		UserID:      iCtx.UserID,
		WorkspaceID: request.GetWorkspaceId(),
		Message:     request.GetMessage(),
	}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &workspaceV1.SendWorkspaceInviteMessageResponse{}, nil
}
