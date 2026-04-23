package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchTaskLogic {
	return &BatchTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchTaskLogic) BatchTask(req *types.BatchTaskReq) (resp *types.BatchTaskResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 0 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	for _, id := range req.Ids {
		task, err := l.svcCtx.TaskModel.FindOne(l.ctx, id)
		if err != nil {
			continue
		}

		if task.UserId != userId {
			continue
		}

		switch req.Action {
		case "complete":
			task.Status = 1
			_ = l.svcCtx.TaskModel.Update(l.ctx, task)
		case "undo":
			task.Status = 0
			_ = l.svcCtx.TaskModel.Update(l.ctx, task)
		case "delete":
			_ = l.svcCtx.TaskModel.Delete(l.ctx, id)
		}
	}

	return &types.BatchTaskResp{}, nil
}
