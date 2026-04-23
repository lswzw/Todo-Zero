// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckRegisterLogic {
	return &CheckRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckRegisterLogic) CheckRegister() (resp *types.CheckRegisterResp, err error) {
	// todo: add your logic here and delete this line

	return
}
