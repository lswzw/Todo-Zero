package admin

import (
	"context"

	"server/internal/pkg/jwtx"
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
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	logs, total, err := l.svcCtx.OperationLogModel.FindList(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.OperationLogItem
	for _, log := range logs {
		list = append(list, types.OperationLogItem{
			Id:         log.Id,
			UserId:     0,
			Username:   log.Username,
			Action:     log.Action,
			TargetType: "",
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

func (l *OperationLogListLogic) checkAdmin() error {
	isAdmin, err := jwtx.GetIsAdminFromCtx(l.ctx)
	if err != nil || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
