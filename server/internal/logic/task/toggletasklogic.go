// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package task

import (
	"context"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleTaskLogic {
	return &ToggleTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleTaskLogic) ToggleTask(req *types.ToggleTaskReq) (resp *types.ToggleTaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
