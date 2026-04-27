package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/scheduler"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BackupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBackupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BackupListLogic {
	return &BackupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BackupListLogic) BackupList() (resp *types.BackupListResp, err error) {
	backups, err := scheduler.ListBackups(l.svcCtx.Config.Database.DataDir)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.BackupItem
	for _, b := range backups {
		list = append(list, types.BackupItem{
			FileName:   b.FileName,
			FileSize:   b.FileSize,
			CreateTime: b.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	if list == nil {
		list = []types.BackupItem{}
	}

	return &types.BackupListResp{List: list}, nil
}
