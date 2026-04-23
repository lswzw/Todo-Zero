// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleUserStatusLogic {
	return &ToggleUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleUserStatusLogic) ToggleUserStatus(req *types.ToggleUserStatusReq) (resp *types.ToggleUserStatusResp, err error) {
	// todo: add your logic here and delete this line

	return
}
