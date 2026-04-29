package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/internal/svc"
	"server/internal/types"
)

func DownloadBackupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadBackupReq
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		fileName := filepath.Base(req.FileName)
		if !strings.HasSuffix(fileName, ".bak") && !strings.HasSuffix(fileName, ".BAK") {
			http.Error(w, "invalid file", http.StatusBadRequest)
			return
		}

		backupDir := filepath.Join(svcCtx.Config.Database.DataDir, "backups")
		absBackupDir, _ := filepath.Abs(backupDir)
		filePath := filepath.Join(backupDir, fileName)
		absPath, err := filepath.Abs(filePath)
		if err != nil || !strings.HasPrefix(absPath, absBackupDir) {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		file, err := os.Open(absPath)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil || info.IsDir() {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
		http.ServeContent(w, r, fileName, info.ModTime(), file)
	}
}
