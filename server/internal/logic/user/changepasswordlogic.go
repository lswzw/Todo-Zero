package user

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	// 1. 获取当前用户
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, xerr.NewCodeError(xerr.OldPasswordError)
	}

	// 4. 更新密码（由 Model 层处理哈希）
	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, userId, req.NewPassword); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ChangePasswordResp{}, nil
}
