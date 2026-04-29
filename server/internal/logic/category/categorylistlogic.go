package category

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListLogic {
	return &CategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryListLogic) CategoryList() (resp *types.CategoryListResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	categories, err := l.svcCtx.CategoryModel.FindAll(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	var list []types.CategoryItem
	for _, c := range categories {
		list = append(list, types.CategoryItem{
			Id:       c.Id,
			Name:     c.Name,
			Color:    c.Color,
			Icon:     c.Icon,
			Sort:     c.Sort,
			IsSystem: c.IsSystem,
		})
	}

	return &types.CategoryListResp{List: list}, nil
}
