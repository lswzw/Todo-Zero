package task

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportTaskLogic {
	return &ExportTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportTaskLogic) ExportTask(req *types.ExportTaskReq, w http.ResponseWriter) error {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}

	var status int64 = -1
	if req.Status == 0 || req.Status == 2 {
		status = req.Status
	}

	var priority int64 = -1
	if req.Priority >= 1 && req.Priority <= 3 {
		priority = req.Priority
	}

	tasks, err := l.svcCtx.TaskModel.FindAllForExport(l.ctx, userId, req.Keyword, status, priority, req.CategoryId)
	if err != nil {
		return xerr.NewCodeError(xerr.ServerCommonError)
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

	// 构建导出数据
	var items []types.TaskItem
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

		items = append(items, types.TaskItem{
			Id:           t.Id,
			Title:        t.Title,
			Content:      content,
			Status:       t.Status,
			Priority:     t.Priority,
			CategoryId:   categoryId,
			CategoryName: categoryName,
			StartTime:    formatNullTimeExport(t.StartTime),
			EndTime:      formatNullTimeExport(t.EndTime),
			Reminder:     formatNullTimeExport(t.Reminder),
			Tags:         t.Tags,
			SortOrder:    t.SortOrder,
			CreateTime:   t.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	switch req.Format {
	case "csv":
		return l.writeCSV(items, w)
	default:
		return l.writeJSON(items, w)
	}
}

func (l *ExportTaskLogic) writeCSV(items []types.TaskItem, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=tasks.csv")

	// BOM for Excel UTF-8 compatibility
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"ID", "标题", "内容", "状态", "优先级", "分类", "开始时间", "截止时间", "提醒时间", "标签", "创建时间"}
	if err := cw.Write(header); err != nil {
		return err
	}

	for _, item := range items {
		statusText := "待办"
		if item.Status == 2 {
			statusText = "已完成"
		}
		priorityText := "普通"
		switch item.Priority {
		case 1:
			priorityText = "重要"
		case 2:
			priorityText = "紧急"
		}

		row := []string{
			fmt.Sprintf("%d", item.Id),
			item.Title,
			item.Content,
			statusText,
			priorityText,
			item.CategoryName,
			item.StartTime,
			item.EndTime,
			item.Reminder,
			item.Tags,
			item.CreateTime,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func (l *ExportTaskLogic) writeJSON(items []types.TaskItem, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=tasks.json")

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return xerr.NewCodeError(xerr.ServerCommonError)
	}

	_, err = w.Write(data)
	return err
}

// formatNullTimeExport 将 sql.NullTime 格式化为字符串，无效值返回空串
func formatNullTimeExport(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format("2006-01-02 15:04")
}
