package task

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskLogic) DeleteTask(req *types.DeleteTaskReq) (resp *types.DeleteTaskResp, err error) {
	userId, ok := l.ctx.Value("userId").(float64)
	if !ok || userId == 0 {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.TaskNotFoundError)
	}

	if task.UserId != int64(userId) {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if err := l.svcCtx.TaskModel.Delete(l.ctx, req.Id); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.DeleteTaskResp{}, nil
}
