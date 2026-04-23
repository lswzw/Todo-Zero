package task

import (
	"context"
	"database/sql"

	"server/internal/model"
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
	userId, ok := l.ctx.Value("userId").(float64)
	if !ok || userId == 0 {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	content := sql.NullString{String: req.Content, Valid: req.Content != ""}
	categoryId := sql.NullInt64{Int64: req.CategoryId, Valid: req.CategoryId != 0}

	result, err := l.svcCtx.TaskModel.Insert(l.ctx, &model.Task{
		Title:      req.Title,
		Content:    content,
		Status:     0,
		Priority:   req.Priority,
		CategoryId: categoryId,
		UserId:     int64(userId),
	})
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	id, _ := result.LastInsertId()
	return &types.CreateTaskResp{Id: id}, nil
}
