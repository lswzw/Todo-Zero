package stat

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatLogic {
	return &StatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatLogic) Stat() (resp *types.StatResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	total, todo, done, _, err := l.svcCtx.TaskModel.CountStats(l.ctx, userId)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var doneRate int64
	if total > 0 {
		doneRate = done * 100 / total
	}

	return &types.StatResp{
		Total:    total,
		Done:     done,
		Todo:     todo,
		DoneRate: doneRate,
	}, nil
}
