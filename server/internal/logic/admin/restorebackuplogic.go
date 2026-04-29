package admin

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"server/internal/pkg/xerr"
	"server/internal/scheduler"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreBackupLogic {
	return &RestoreBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreBackupLogic) RestoreBackup(req *types.RestoreBackupReq) (resp *types.RestoreBackupResp, err error) {
	dataDir := l.svcCtx.Config.Database.DataDir
	dbFile := l.svcCtx.Config.Database.DBFile

	// Record the pre-restore safety backup file name before restore
	backupDir := filepath.Join(dataDir, "backups")
	preRestoreName := fmt.Sprintf("%s_prerestore_%s.bak",
		strings.TrimSuffix(dbFile, filepath.Ext(dbFile)),
		time.Now().Format("20060102_150405"))

	if err := scheduler.RestoreBackup(l.svcCtx.DB, dataDir, dbFile, req.FileName); err != nil {
		l.Errorf("Restore backup failed: %v", err)
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// Clean up pre-restore safety backups beyond max count
	dbBackupMaxCount := scheduler.GetConfigInt(l.svcCtx, "db_backup_max_count", 7)
	if dbBackupMaxCount > 0 {
		scheduler.CleanOldBackups(backupDir, int(dbBackupMaxCount)+2) // +2 to keep safety backups
	}

	return &types.RestoreBackupResp{
		PreRestoreBackup: preRestoreName,
	}, nil
}
