package task

import (
	"context"
	"database/sql"

	"server/internal/pkg/jwtx"
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
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.TaskNotFoundError)
	}

	// 只能操作自己的任务
	if task.UserId != userId {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	// 使用指针类型区分"未提供"(nil)和"清空"(非nil零值)
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Content != nil {
		task.Content = sql.NullString{String: *req.Content, Valid: *req.Content != ""}
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.CategoryId != nil {
		task.CategoryId = sql.NullInt64{Int64: *req.CategoryId, Valid: *req.CategoryId != 0}
	}
	if req.StartTime != nil {
		task.StartTime = parseNullTime(*req.StartTime)
	}
	if req.EndTime != nil {
		task.EndTime = parseNullTime(*req.EndTime)
	}
	if req.Reminder != nil {
		task.Reminder = parseNullTime(*req.Reminder)
	}
	if req.Tags != nil {
		task.Tags = *req.Tags
	}

	if err := l.svcCtx.TaskModel.Update(l.ctx, task); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.UpdateTaskResp{}, nil
}
