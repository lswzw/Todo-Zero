package admin

import (
	"context"

	"server/internal/pkg/jwtx"
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
	// 不能删除自己
	userId, _ := jwtx.GetUserIdFromCtx(l.ctx)
	if userId == req.Id {
		return nil, xerr.NewCodeErrFromMsg("不能删除自己")
	}

	// 不能删除其他管理员
	targetUser, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}
	if targetUser.Role == 1 {
		return nil, xerr.NewCodeErrFromMsg("不能删除管理员账户")
	}

	if err := l.svcCtx.UserModel.Delete(l.ctx, req.Id); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.DeleteUserResp{}, nil
}
