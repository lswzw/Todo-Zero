package category

import (
	"context"
	"database/sql"

	"server/internal/model"
	"server/internal/pkg/jwtx"
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
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.CategoryModel.Insert(l.ctx, &model.Category{
		Name:   req.Name,
		Color:  "#1890ff",
		UserId: sql.NullInt64{Int64: userId, Valid: true},
		Sort:   0,
	})
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	id, _ := result.LastInsertId()
	return &types.CreateCategoryResp{Id: id}, nil
}
