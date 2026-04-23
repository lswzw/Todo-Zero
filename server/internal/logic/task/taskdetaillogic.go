// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
