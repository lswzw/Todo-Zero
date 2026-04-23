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

	tasks, _, err := l.svcCtx.TaskModel.FindList(l.ctx, userId, 0, 0, 0, "", 1, 9999)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var total, done int64
	total = int64(len(tasks))
	for _, t := range tasks {
		if t.Status == 1 {
			done++
		}
	}

	todo := total - done
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
