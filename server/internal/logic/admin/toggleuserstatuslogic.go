package admin

import (
	"context"

	"server/internal/pkg/xerr"
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
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	// 切换状态
	if user.Status == 1 {
		user.Status = 0
	} else {
		user.Status = 1
	}

	if err := l.svcCtx.UserModel.Update(l.ctx, user); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ToggleUserStatusResp{}, nil
}

func (l *ToggleUserStatusLogic) checkAdmin() error {
	isAdmin, ok := l.ctx.Value("isAdmin").(float64)
	if !ok || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
