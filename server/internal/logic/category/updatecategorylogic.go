package category

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryLogic {
	return &UpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(req *types.UpdateCategoryReq) (resp *types.UpdateCategoryResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	cat, err := l.svcCtx.CategoryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.CategoryNotFoundError)
	}

	// 只能修改自己的分类，系统分类不可修改
	if cat.IsSystem == 1 || (cat.UserId.Valid && cat.UserId.Int64 != userId) {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	// 更新字段
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Color != nil {
		cat.Color = *req.Color
	}
	if req.Icon != nil {
		cat.Icon = *req.Icon
	}
	if req.Sort != nil {
		cat.Sort = *req.Sort
	}

	if err := l.svcCtx.CategoryModel.Update(l.ctx, cat); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.UpdateCategoryResp{}, nil
}
