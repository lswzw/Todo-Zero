package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleTaskLogic {
	return &ToggleTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleTaskLogic) ToggleTask(req *types.ToggleTaskReq) (resp *types.ToggleTaskResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.TaskNotFoundError)
	}

	if task.UserId != userId {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	// Atomic status toggle: 0=待办 ↔ 2=已完成
	newStatus := int64(0)
	if task.Status == 0 {
		newStatus = 2
	}
	if err := l.svcCtx.TaskModel.UpdateStatus(l.ctx, req.Id, newStatus); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ToggleTaskResp{}, nil
}
