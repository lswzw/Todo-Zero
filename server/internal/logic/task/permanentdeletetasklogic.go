package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermanentDeleteTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermanentDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermanentDeleteTaskLogic {
	return &PermanentDeleteTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermanentDeleteTaskLogic) PermanentDeleteTask(req *types.PermanentDeleteTaskReq) (resp *types.PermanentDeleteTaskResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	task, err := l.svcCtx.TaskModel.FindOneIncludeDeleted(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.TaskNotFoundError)
	}

	if task.UserId != userId {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if task.IsDeleted == 0 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	if err := l.svcCtx.TaskModel.PermanentDelete(l.ctx, req.Id); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.PermanentDeleteTaskResp{}, nil
}
