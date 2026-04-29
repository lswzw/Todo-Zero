package admin

import (
	"context"

	"server/internal/model"
	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.ResetPasswordResp, err error) {
	currentUserId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	isAdmin, err := jwtx.GetIsAdminFromCtx(l.ctx)
	if err != nil || isAdmin != 1 {
		return nil, xerr.NewCodeError(xerr.AdminRequired)
	}

	if currentUserId == req.Id {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	targetUser, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.UserNotFoundError)
		}
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	if targetUser.Role == 1 {
		return nil, xerr.NewCodeError(xerr.AdminRequired)
	}

	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, req.Id, req.NewPassword); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ResetPasswordResp{}, nil
}
