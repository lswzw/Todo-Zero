package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskListLogic) TaskList(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	tasks, total, err := l.svcCtx.TaskModel.FindList(l.ctx, userId, req.Status, req.CategoryId, req.Priority, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.TaskItem
	for _, t := range tasks {
		categoryName := "未分类"
		if t.CategoryId.Valid && t.CategoryId.Int64 > 0 {
			category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, t.CategoryId.Int64)
			if err == nil {
				categoryName = category.Name
			}
		}

		list = append(list, types.TaskItem{
			Id:           t.Id,
			Title:        t.Title,
			Content:      t.Content.String,
			Status:       t.Status,
			Priority:     t.Priority,
			CategoryId:   t.CategoryId.Int64,
			CategoryName: categoryName,
			CreateTime:   t.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	return &types.TaskListResp{
		Total: total,
		List:  list,
	}, nil
}
