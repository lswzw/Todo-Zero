package task

import (
	"context"

	"server/internal/model"
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

	switch req.Action {
	case "complete", "undo", "delete", "restore":
	default:
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	if len(req.Ids) == 0 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	var failedIds []int64
	for _, id := range req.Ids {
		// restore 需要查找已删除的任务，其他操作只查找未删除的
		var task *model.Task
		var err error
		if req.Action == "restore" {
			task, err = l.svcCtx.TaskModel.FindOneIncludeDeleted(l.ctx, id)
		} else {
			task, err = l.svcCtx.TaskModel.FindOne(l.ctx, id)
		}
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if task.UserId != userId {
			failedIds = append(failedIds, id)
			continue
		}

		switch req.Action {
		case "complete":
			if err := l.svcCtx.TaskModel.UpdateStatus(l.ctx, id, 2); err != nil {
				failedIds = append(failedIds, id)
			}
		case "undo":
			if err := l.svcCtx.TaskModel.UpdateStatus(l.ctx, id, 0); err != nil {
				failedIds = append(failedIds, id)
			}
		case "delete":
			if err := l.svcCtx.TaskModel.Delete(l.ctx, id); err != nil {
				failedIds = append(failedIds, id)
			}
		case "restore":
			if err := l.svcCtx.TaskModel.Restore(l.ctx, id, userId); err != nil {
				failedIds = append(failedIds, id)
			}
		}
	}

	if len(failedIds) > 0 {
		logx.Errorf("[BatchTask] failed ids: %v, action: %s", failedIds, req.Action)
	}

	return &types.BatchTaskResp{}, nil
}
