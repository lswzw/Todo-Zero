package task

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrashListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrashListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrashListLogic {
	return &TrashListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrashListLogic) TrashList(req *types.TrashListReq) (resp *types.TrashListResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	tasks, total, err := l.svcCtx.TaskModel.FindDeletedList(l.ctx, userId, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 批量收集 categoryId
	categoryIds := make(map[int64]bool)
	for _, t := range tasks {
		if t.CategoryId.Valid && t.CategoryId.Int64 > 0 {
			categoryIds[t.CategoryId.Int64] = true
		}
	}
	categoryMap := make(map[int64]string)
	for cid := range categoryIds {
		category, err := l.svcCtx.CategoryModel.FindOne(l.ctx, cid)
		if err == nil {
			categoryMap[cid] = category.Name
		}
	}

	var list []types.TrashItem
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

		list = append(list, types.TrashItem{
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
			UpdateTime:   t.UpdateTime.Format("2006-01-02 15:04"),
		})
	}

	return &types.TrashListResp{
		Total: total,
		List:  list,
	}, nil
}
