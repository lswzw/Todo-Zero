package admin

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, req.Id, string(hashedPassword)); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ResetPasswordResp{}, nil
}

func (l *ResetPasswordLogic) checkAdmin() error {
	isAdmin, err := jwtx.GetIsAdminFromCtx(l.ctx)
	if err != nil || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
