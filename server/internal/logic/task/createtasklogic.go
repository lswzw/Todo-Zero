package task

import (
	"context"
	"database/sql"
	"time"

	"server/internal/model"
	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.CreateTaskResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	content := sql.NullString{String: req.Content, Valid: req.Content != ""}
	categoryId := sql.NullInt64{Int64: req.CategoryId, Valid: req.CategoryId != 0}
	startTime := parseNullTime(req.StartTime)
	endTime := parseNullTime(req.EndTime)
	reminder := parseNullTime(req.Reminder)

	result, err := l.svcCtx.TaskModel.Insert(l.ctx, &model.Task{
		Title:      req.Title,
		Content:    content,
		Status:     0,
		Priority:   req.Priority,
		CategoryId: categoryId,
		UserId:     userId,
		StartTime:  startTime,
		EndTime:    endTime,
		Reminder:   reminder,
		Tags:       req.Tags,
	})
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	id, _ := result.LastInsertId()
	return &types.CreateTaskResp{Id: id}, nil
}

// parseNullTime 解析时间字符串为 sql.NullTime，空字符串返回无效值
func parseNullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{Valid: false}
	}
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// formatNullTime 将 sql.NullTime 格式化为字符串，无效值返回空串
func formatNullTime(nt sql.NullTime) string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format("2006-01-02 15:04")
}
