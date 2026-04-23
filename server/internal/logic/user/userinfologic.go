package user

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResp, err error) {
	userId, err := l.getJwtUserId()
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.UserNotFoundError)
	}

	return &types.UserInfoResp{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.Role,
		Status:   user.Status,
	}, nil
}

// getJwtUserId 从 context 获取 JWT 中的 userId
func (l *UserInfoLogic) getJwtUserId() (int64, error) {
	return jwtx.GetUserIdFromCtx(l.ctx)
}
