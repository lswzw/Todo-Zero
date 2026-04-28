package tag

import (
	"context"

	"server/internal/model"
	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagReq) (resp *types.CreateTagResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if _, err := l.svcCtx.TagModel.FindByName(l.ctx, userId, req.Name); err == nil {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	result, err := l.svcCtx.TagModel.Insert(l.ctx, &model.Tag{
		Name:     req.Name,
		Color:    req.Color,
		UserId:   userId,
		IsSystem: 0,
	})
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &types.CreateTagResp{Id: id}, nil
}
