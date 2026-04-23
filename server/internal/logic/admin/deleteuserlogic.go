package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	// 不能删除自己
	userId, _ := l.ctx.Value("userId").(float64)
	if int64(userId) == req.Id {
		return nil, xerr.NewCodeErrFromMsg("不能删除自己")
	}

	if err := l.svcCtx.UserModel.Delete(l.ctx, req.Id); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.DeleteUserResp{}, nil
}

func (l *DeleteUserLogic) checkAdmin() error {
	isAdmin, ok := l.ctx.Value("isAdmin").(float64)
	if !ok || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
