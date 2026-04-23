package task

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskLogic) UpdateTask(req *types.UpdateTaskReq) (resp *types.UpdateTaskResp, err error) {
	userId, ok := l.ctx.Value("userId").(float64)
	if !ok || userId == 0 {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.TaskNotFoundError)
	}

	// 只能操作自己的任务
	if task.UserId != int64(userId) {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	// 更新字段
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Content != "" {
		task.Content.String = req.Content
		task.Content.Valid = true
	}
	if req.Priority != 0 {
		task.Priority = req.Priority
	}
	if req.CategoryId != 0 {
		task.CategoryId.Int64 = req.CategoryId
		task.CategoryId.Valid = true
	}

	if err := l.svcCtx.TaskModel.Update(l.ctx, task); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.UpdateTaskResp{}, nil
}
