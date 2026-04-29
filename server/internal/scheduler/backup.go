package scheduler

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"server/internal/svc"
)

const backupCheckInterval = 6 * time.Hour

var allowedTables = map[string]bool{
	"tasks":          true,
	"operation_logs": true,
	"login_log":      true,
	"categories":     true,
	"system_configs": true,
	"users":          true,
}

func isValidTableName(table string) bool {
	return allowedTables[table]
}

// StartBackupScheduler starts a background goroutine that periodically backs up the SQLite database.
func StartBackupScheduler(svcCtx *svc.ServiceContext) {
	go func() {
		// Delay first run to let the server finish startup
		time.Sleep(1 * time.Minute)
		runBackup(svcCtx)

		ticker := time.NewTicker(backupCheckInterval)
		defer ticker.Stop()
		for range ticker.C {
			runBackup(svcCtx)
		}
	}()
	fmt.Println("[Scheduler] Backup scheduler started (check interval: 6h)")
}

// runBackup checks if a backup is needed and performs it.
func runBackup(svcCtx *svc.ServiceContext) {
	dbBackupEnabled := getConfigInt(svcCtx, "db_backup_enabled", 0)
	if dbBackupEnabled == 0 {
		return
	}

	dbBackupIntervalHours := getConfigInt(svcCtx, "db_backup_interval_hours", 24)
	dbBackupMaxCount := getConfigInt(svcCtx, "db_backup_max_count", 7)

	dataDir := svcCtx.Config.Database.DataDir
	dbFile := svcCtx.Config.Database.DBFile
	backupDir := filepath.Join(dataDir, "backups")

	// Ensure backup directory exists with restrictive permissions (owner read/write only)
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		fmt.Printf("[Scheduler-Backup] Failed to create backup directory: %v\n", err)
		return
	}

	// Check if backup is needed (compare with latest backup time)
	if !needsBackup(backupDir, dbBackupIntervalHours) {
		return
	}

	// Perform backup using VACUUM INTO on the existing connection
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("%s_%s.bak", strings.TrimSuffix(dbFile, filepath.Ext(dbFile)), timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	if err := PerformBackup(svcCtx.DB, backupPath); err != nil {
		fmt.Printf("[Scheduler-Backup] Failed to backup database: %v\n", err)
		return
	}

	fmt.Printf("[Scheduler-Backup] Database backed up to %s\n", backupPath)

	// Clean up old backups
	if dbBackupMaxCount > 0 {
		cleanOldBackups(backupDir, int(dbBackupMaxCount))
	}
}

// needsBackup checks if enough time has passed since the last backup.
func needsBackup(backupDir string, intervalHours int64) bool {
	files, err := ioutil.ReadDir(backupDir)
	if err != nil || len(files) == 0 {
		return true // No backups yet
	}

	// Find the most recent backup
	var latestTime time.Time
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".bak") && f.ModTime().After(latestTime) {
			latestTime = f.ModTime()
		}
	}

	elapsed := time.Since(latestTime)
	return elapsed >= time.Duration(intervalHours)*time.Hour
}

// PerformBackup creates a consistent backup of the SQLite database using VACUUM INTO.
// This is exported for use by the manual backup API.
func PerformBackup(db *sql.DB, backupPath string) error {
	// Ensure backup directory exists with restrictive permissions
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	vacuumSQL := fmt.Sprintf("VACUUM INTO '%s'", backupPath)

	// Remove existing file if any (VACUUM INTO fails on existing files)
	os.Remove(backupPath)

	if _, err := db.Exec(vacuumSQL); err != nil {
		// Clean up failed backup file
		os.Remove(backupPath)
		return fmt.Errorf("VACUUM INTO failed: %w", err)
	}

	// Set restrictive permissions on backup file (owner read/write only)
	if err := os.Chmod(backupPath, 0600); err != nil {
		return fmt.Errorf("failed to set backup file permissions: %w", err)
	}

	// Verify backup file exists and has content
	info, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("backup file not found after VACUUM INTO: %w", err)
	}
	if info.Size() == 0 {
		os.Remove(backupPath)
		return fmt.Errorf("backup file is empty")
	}

	return nil
}

// BackupInfo represents metadata about a backup file.
type BackupInfo struct {
	FileName   string    `json:"fileName"`
	FileSize   int64     `json:"fileSize"`
	CreateTime time.Time `json:"createTime"`
}

// ListBackups returns a list of backup files sorted by creation time (newest first).
func ListBackups(dataDir string) ([]BackupInfo, error) {
	backupDir := filepath.Join(dataDir, "backups")
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var backups []BackupInfo
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".bak") {
			continue
		}
		backups = append(backups, BackupInfo{
			FileName:   f.Name(),
			FileSize:   f.Size(),
			CreateTime: f.ModTime(),
		})
	}

	// Sort by create time descending (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreateTime.After(backups[j].CreateTime)
	})

	return backups, nil
}

// RestoreBackup replaces the current database data with data from a backup file.
// It uses ATTACH DATABASE to copy data within a transaction.
// A safety backup of the current database is created before restoring.
func RestoreBackup(db *sql.DB, dataDir string, dbFile string, backupFileName string) error {
	// Validate backup file name
	if !strings.HasSuffix(backupFileName, ".bak") {
		return fmt.Errorf("invalid backup file")
	}

	backupDir := filepath.Join(dataDir, "backups")
	backupPath := filepath.Join(backupDir, backupFileName)

	// Security: prevent path traversal
	absBackupDir, _ := filepath.Abs(backupDir)
	absBackupPath, err := filepath.Abs(backupPath)
	if err != nil || !strings.HasPrefix(absBackupPath, absBackupDir) {
		return fmt.Errorf("invalid backup file path")
	}

	// Check file exists
	info, err := os.Stat(absBackupPath)
	if err != nil || info.IsDir() {
		return fmt.Errorf("backup file not found")
	}
	if info.Size() == 0 {
		return fmt.Errorf("backup file is empty")
	}

	// Create a safety backup before restoring
	safetyBackupPath := filepath.Join(backupDir,
		fmt.Sprintf("%s_prerestore_%s.bak",
			strings.TrimSuffix(dbFile, filepath.Ext(dbFile)),
			time.Now().Format("20060102_150405")))
	if err := PerformBackup(db, safetyBackupPath); err != nil {
		return fmt.Errorf("failed to create safety backup before restore: %w", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// Attach backup database
	_, err = tx.Exec(fmt.Sprintf("ATTACH DATABASE '%s' AS backup_db", absBackupPath))
	if err != nil {
		return fmt.Errorf("failed to attach backup database: %w", err)
	}

	// Verify backup has required tables
	var tableCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM backup_db.sqlite_master WHERE type='table' AND name IN ('users','categories','tasks','system_configs')").Scan(&tableCount)
	if err != nil || tableCount < 4 {
		tx.Rollback()
		committed = true
		db.Exec("DETACH DATABASE backup_db")
		return fmt.Errorf("backup file does not contain required tables")
	}

	// Delete data in order (children first to respect foreign keys)
	deleteOrder := []string{"tasks", "operation_logs", "login_log", "categories", "system_configs", "users"}
	for _, table := range deleteOrder {
		if !isValidTableName(table) {
			tx.Rollback()
			committed = true
			db.Exec("DETACH DATABASE backup_db")
			return fmt.Errorf("invalid table name: %s", table)
		}
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}

	// Insert data in order (parents first to satisfy foreign keys)
	insertOrder := []string{"users", "categories", "system_configs", "tasks", "operation_logs", "login_log"}
	for _, table := range insertOrder {
		if !isValidTableName(table) {
			tx.Rollback()
			committed = true
			db.Exec("DETACH DATABASE backup_db")
			return fmt.Errorf("invalid table name: %s", table)
		}
		if _, err := tx.Exec(fmt.Sprintf("INSERT INTO %s SELECT * FROM backup_db.%s", table, table)); err != nil {
			return fmt.Errorf("failed to restore %s: %w", table, err)
		}
	}

	// Commit first, then detach (DETACH inside transaction causes "database is locked")
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit restore: %w", err)
	}
	committed = true

	// Detach backup database after commit
	if _, err := db.Exec("DETACH DATABASE backup_db"); err != nil {
		fmt.Printf("[Scheduler-Backup] Warning: failed to detach backup database: %v\n", err)
	}

	return nil
}

// GetConfigInt reads an integer config value from system_configs table.
func GetConfigInt(svcCtx *svc.ServiceContext, key string, defaultVal int64) int64 {
	return getConfigInt(svcCtx, key, defaultVal)
}

// CleanOldBackups removes old backup files, keeping only the most recent maxCount.
func CleanOldBackups(backupDir string, maxCount int) {
	cleanOldBackups(backupDir, maxCount)
}

// cleanOldBackups removes old backup files, keeping only the most recent maxCount.
func cleanOldBackups(backupDir string, maxCount int) {
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return
	}

	// Filter .bak files
	var bakFiles []os.FileInfo
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".bak") {
			bakFiles = append(bakFiles, f)
		}
	}

	// Sort by modification time descending (newest first)
	sort.Slice(bakFiles, func(i, j int) bool {
		return bakFiles[i].ModTime().After(bakFiles[j].ModTime())
	})

	// Remove excess files
	for i := maxCount; i < len(bakFiles); i++ {
		path := filepath.Join(backupDir, bakFiles[i].Name())
		if err := os.Remove(path); err != nil {
			fmt.Printf("[Scheduler-Backup] Failed to remove old backup %s: %v\n", bakFiles[i].Name(), err)
		} else {
			fmt.Printf("[Scheduler-Backup] Removed old backup: %s\n", bakFiles[i].Name())
		}
	}
}
