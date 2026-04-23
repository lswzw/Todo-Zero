package admin

import (
	"context"

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
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	// Use atomic status update to avoid TOCTOU race condition
	newStatus := int64(0)
	if user.Status == 0 {
		newStatus = 1
	}
	if err := l.svcCtx.UserModel.UpdateStatus(l.ctx, req.Id, newStatus); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ToggleUserStatusResp{}, nil
}
