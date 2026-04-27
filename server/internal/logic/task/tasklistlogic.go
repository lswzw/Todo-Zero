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

	// Status: -1=全部, 0=待办, 2=已完成
	// Priority: -1=全部, 1=重要, 2=紧急, 3=普通
	// 用 default=-1 区分"未传"和"传了0"，解决 go-zero optional int 零值问题
	var status int64 = -1
	if req.Status == 0 || req.Status == 2 {
		status = req.Status
	}

	var priority int64 = -1
	if req.Priority >= 1 && req.Priority <= 3 {
		priority = req.Priority
	}

	tasks, total, err := l.svcCtx.TaskModel.FindList(l.ctx, userId, req.Keyword, status, priority, req.CategoryId, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 批量收集 categoryId，一次查出所有分类名称
	categoryIds := make(map[int64]bool)
	for _, t := range tasks {
		if t.CategoryId.Valid && t.CategoryId.Int64 > 0 {
			categoryIds[t.CategoryId.Int64] = true
		}
	}
	categoryMap := make(map[int64]string)
	if len(categoryIds) > 0 {
		for cid := range categoryIds {
			category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, cid)
			if err == nil {
				categoryMap[cid] = category.Name
			}
		}
	}

	var list []types.TaskItem
	for _, t := range tasks {
		categoryName := "未分类"
		categoryId := int64(0)
		if t.CategoryId.Valid && t.CategoryId.Int64 > 0 {
			categoryId = t.CategoryId.Int64
			if name, ok := categoryMap[categoryId]; ok {
				categoryName = name
			}
		}

		content := ""
		if t.Content.Valid {
			content = t.Content.String
		}

		list = append(list, types.TaskItem{
			Id:           t.Id,
			Title:        t.Title,
			Content:      content,
			Status:       t.Status,
			Priority:     t.Priority,
			CategoryId:   categoryId,
			CategoryName: categoryName,
			StartTime:    formatNullTime(t.StartTime),
			EndTime:      formatNullTime(t.EndTime),
			Reminder:     formatNullTime(t.Reminder),
			Tags:         t.Tags,
			CreateTime:   t.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	return &types.TaskListResp{
		Total: total,
		List:  list,
	}, nil
}
