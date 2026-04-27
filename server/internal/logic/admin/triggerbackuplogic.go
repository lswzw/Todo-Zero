package admin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/internal/pkg/xerr"
	"server/internal/scheduler"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerBackupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerBackupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerBackupLogic {
	return &TriggerBackupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerBackupLogic) TriggerBackup(req *types.TriggerBackupReq) (resp *types.TriggerBackupResp, err error) {
	dataDir := l.svcCtx.Config.Database.DataDir
	dbFile := l.svcCtx.Config.Database.DBFile
	backupDir := filepath.Join(dataDir, "backups")

	// Ensure backup directory exists
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("%s_%s_manual.bak", strings.TrimSuffix(dbFile, filepath.Ext(dbFile)), timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	if err := scheduler.PerformBackup(l.svcCtx.DB, backupPath); err != nil {
		l.Errorf("Manual backup failed: %v", err)
		return nil, xerr.NewCodeErrFromMsg("备份失败: " + err.Error())
	}

	// Get file size
	var fileSize int64
	if info, statErr := os.Stat(backupPath); statErr == nil {
		fileSize = info.Size()
	}

	return &types.TriggerBackupResp{
		FileName: backupFileName,
		FileSize: fileSize,
	}, nil
}
