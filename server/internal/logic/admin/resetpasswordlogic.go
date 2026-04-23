package admin

import (
	"context"

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

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	user.Password = string(hashedPassword)
	if err := l.svcCtx.UserModel.Update(l.ctx, user); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.ResetPasswordResp{}, nil
}

func (l *ResetPasswordLogic) checkAdmin() error {
	isAdmin, ok := l.ctx.Value("isAdmin").(float64)
	if !ok || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
