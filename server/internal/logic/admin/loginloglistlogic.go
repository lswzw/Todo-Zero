package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogListLogic {
	return &LoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogListLogic) LoginLogList(req *types.LoginLogReq) (resp *types.LoginLogResp, err error) {
	logs, total, err := l.svcCtx.LoginLogModel.FindList(l.ctx, req.Username, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.LoginLogItem
	for _, log := range logs {
		userId := int64(0)
		if log.UserId.Valid {
			userId = log.UserId.Int64
		}
		list = append(list, types.LoginLogItem{
			Id:         log.Id,
			UserId:     userId,
			Username:   log.Username,
			Ip:         log.Ip,
			Status:     log.Status,
			Remark:     log.Remark,
			CreateTime: log.CreateTime.Format("2006-01-02 15:04"),
		})
	}

	return &types.LoginLogResp{
		Total: total,
		List:  list,
	}, nil
}
