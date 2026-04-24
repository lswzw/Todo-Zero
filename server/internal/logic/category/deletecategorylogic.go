package category

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(req *types.DeleteCategoryReq) (resp *types.DeleteCategoryResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	cat, err := l.svcCtx.CategoryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.CategoryNotFoundError)
	}

	// 只能删除自己的分类，系统分类不可删除
	if cat.IsSystem == 1 || (cat.UserId.Valid && cat.UserId.Int64 != userId) {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if err := l.svcCtx.CategoryModel.Delete(l.ctx, req.Id); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.DeleteCategoryResp{}, nil
}
