package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (resp *types.UserListResp, err error) {
	users, total, err := l.svcCtx.UserModel.FindList(l.ctx, req.Keyword, -1, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.UserListItem
	for _, u := range users {
		list = append(list, types.UserListItem{
			Id:         u.Id,
			Username:   u.Username,
			IsAdmin:    u.Role,
			Status:     u.Status,
			CreateTime: u.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	return &types.UserListResp{
		Total: total,
		List:  list,
	}, nil
}
