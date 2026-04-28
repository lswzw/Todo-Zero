package tag

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTagLogic {
	return &UpdateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTagLogic) UpdateTag(req *types.UpdateTagReq) (resp *types.UpdateTagResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if tag.UserId != userId {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if tag.IsSystem == 1 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	if req.Name != nil && *req.Name != "" {
		existing, _ := l.svcCtx.TagModel.FindByName(l.ctx, userId, *req.Name)
		if existing != nil && existing.Id != req.Id {
			return nil, xerr.NewCodeError(xerr.RequestParamError)
		}
		tag.Name = *req.Name
	}

	if req.Color != nil && *req.Color != "" {
		tag.Color = *req.Color
	}

	err = l.svcCtx.TagModel.Update(l.ctx, tag)
	if err != nil {
		return nil, err
	}

	return &types.UpdateTagResp{Success: true}, nil
}
