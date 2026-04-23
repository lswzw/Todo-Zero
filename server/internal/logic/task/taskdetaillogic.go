package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDetailLogic {
	return &TaskDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskDetailLogic) TaskDetail(req *types.TaskDetailReq) (resp *types.TaskDetailResp, err error) {
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

	// 获取分类名称
	categoryName := "未分类"
	if task.CategoryId.Valid && task.CategoryId.Int64 > 0 {
		category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, task.CategoryId.Int64)
		if err == nil {
			categoryName = category.Name
		}
	}

	return &types.TaskDetailResp{
		Id:           task.Id,
		Title:        task.Title,
		Content:      task.Content.String,
		Status:       task.Status,
		Priority:     task.Priority,
		CategoryId:   task.CategoryId.Int64,
		CategoryName: categoryName,
		CreateTime:   task.CreateTime.Format("2006-01-02 15:04"),
		UpdateTime:   task.UpdateTime.Format("2006-01-02 15:04"),
	}, nil
}
