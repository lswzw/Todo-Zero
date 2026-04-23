package category

import (
	"context"
	"database/sql"

	"server/internal/model"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CreateCategoryResp, err error) {
	userId, ok := l.ctx.Value("userId").(float64)
	if !ok || userId == 0 {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	result, err := l.svcCtx.CategoryModel.Insert(l.ctx, &model.Category{
		Name:      req.Name,
		UserId:    sql.NullInt64{Int64: int64(userId), Valid: true},
		SortOrder: 0,
	})
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	id, _ := result.LastInsertId()
	return &types.CreateCategoryResp{Id: id}, nil
}
