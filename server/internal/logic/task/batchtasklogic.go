// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchTaskLogic {
	return &BatchTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchTaskLogic) BatchTask(req *types.BatchTaskReq) (resp *types.BatchTaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
