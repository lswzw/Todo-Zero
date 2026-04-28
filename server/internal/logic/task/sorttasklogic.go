package task

import (
	"context"

	"server/internal/model"
	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SortTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortTaskLogic {
	return &SortTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortTaskLogic) SortTask(req *types.SortTaskReq) (resp *types.SortTaskResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Orders) == 0 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	// 校验所有任务都属于当前用户
	for _, item := range req.Orders {
		task, err := l.svcCtx.TaskModel.FindOne(l.ctx, item.Id)
		if err != nil {
			return nil, xerr.NewCodeError(xerr.RequestParamError)
		}
		if task.UserId != userId {
			return nil, xerr.NewCodeError(xerr.NoPermission)
		}
	}

	orders := make([]model.SortOrderItem, 0, len(req.Orders))
	for _, item := range req.Orders {
		orders = append(orders, model.SortOrderItem{
			Id:        item.Id,
			SortOrder: item.SortOrder,
		})
	}

	if err := l.svcCtx.TaskModel.UpdateSortOrder(l.ctx, userId, orders); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.SortTaskResp{}, nil
}
