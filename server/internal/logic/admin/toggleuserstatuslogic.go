package admin

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleUserStatusLogic {
	return &ToggleUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleUserStatusLogic) ToggleUserStatus(req *types.ToggleUserStatusReq) (resp *types.ToggleUserStatusResp, err error) {
	currentUserId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	isAdmin, err := jwtx.GetIsAdminFromCtx(l.ctx)
	if err != nil || isAdmin != 1 {
		return nil, xerr.NewCodeError(xerr.AdminRequired)
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	if currentUserId == req.Id {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	if user.Role == 1 {
		return nil, xerr.NewCodeError(xerr.AdminRequired)
	}

	newStatus := int64(0)
	if user.Status == 0 {
		newStatus = 1
	}
	if err := l.svcCtx.UserModel.UpdateStatus(l.ctx, req.Id, newStatus); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ToggleUserStatusResp{}, nil
}
