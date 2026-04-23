package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogListLogic {
	return &OperationLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationLogListLogic) OperationLogList(req *types.OperationLogReq) (resp *types.OperationLogResp, err error) {
	logs, total, err := l.svcCtx.OperationLogModel.FindList(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.OperationLogItem
	for _, log := range logs {
		userId := int64(0)
		if log.UserId.Valid {
			userId = log.UserId.Int64
		}
		list = append(list, types.OperationLogItem{
			Id:         log.Id,
			UserId:     userId,
			Username:   log.Username,
			Action:     log.Action,
			TargetType: log.Module,
			TargetId:   0,
			Detail:     log.Params,
			Ip:         log.Ip,
			CreateTime: log.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return &types.OperationLogResp{
		Total: total,
		List:  list,
	}, nil
}
