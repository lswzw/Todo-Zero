package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"server/internal/svc"
)

const checkInterval = 1 * time.Hour

// StartCleanupScheduler starts a background goroutine that periodically cleans up old data.
// Cleanup rules are read from system_configs so changes take effect without restart.
func StartCleanupScheduler(svcCtx *svc.ServiceContext) {
	go func() {
		// Delay first run to let the server finish startup
		time.Sleep(30 * time.Second)
		runCleanup(svcCtx)

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		for range ticker.C {
			runCleanup(svcCtx)
		}
	}()
	fmt.Println("[Scheduler] Cleanup scheduler started (interval: 1h)")
}

func runCleanup(svcCtx *svc.ServiceContext) {
	ctx := context.Background()

	// Read configs from DB (hot-reload on each run)
	taskAutoDeleteDays := getConfigInt(svcCtx, "task_auto_delete_days", 0)
	taskTrashRetentionDays := getConfigInt(svcCtx, "task_trash_retention_days", 30)
	logAutoDeleteDays := getConfigInt(svcCtx, "log_auto_delete_days", 0)

	now := time.Now()

	// 1. Hard-delete completed tasks older than task_auto_delete_days
	if taskAutoDeleteDays > 0 {
		before := now.AddDate(0, 0, -int(taskAutoDeleteDays))
		count, err := svcCtx.TaskModel.HardDeleteCompletedBefore(ctx, before)
		if err != nil {
			fmt.Printf("[Scheduler] Failed to cleanup completed tasks: %v\n", err)
		} else if count > 0 {
			fmt.Printf("[Scheduler] Cleaned up %d completed tasks older than %d days\n", count, taskAutoDeleteDays)
		}
	}

	// 2. Hard-delete soft-deleted tasks older than task_trash_retention_days
	if taskTrashRetentionDays > 0 {
		before := now.AddDate(0, 0, -int(taskTrashRetentionDays))
		count, err := svcCtx.TaskModel.HardDeleteSoftDeletedBefore(ctx, before)
		if err != nil {
			fmt.Printf("[Scheduler] Failed to cleanup soft-deleted tasks: %v\n", err)
		} else if count > 0 {
			fmt.Printf("[Scheduler] Cleaned up %d soft-deleted tasks older than %d days\n", count, taskTrashRetentionDays)
		}
	}

	// 3. Delete old operation logs
	if logAutoDeleteDays > 0 {
		before := now.AddDate(0, 0, -int(logAutoDeleteDays))
		count, err := svcCtx.OperationLogModel.DeleteOlderThan(ctx, before)
		if err != nil {
			fmt.Printf("[Scheduler] Failed to cleanup operation logs: %v\n", err)
		} else if count > 0 {
			fmt.Printf("[Scheduler] Cleaned up %d operation logs older than %d days\n", count, logAutoDeleteDays)
		}
	}

	// 4. Delete old login logs
	if logAutoDeleteDays > 0 {
		before := now.AddDate(0, 0, -int(logAutoDeleteDays))
		count, err := svcCtx.LoginLogModel.DeleteOlderThan(ctx, before)
		if err != nil {
			fmt.Printf("[Scheduler] Failed to cleanup login logs: %v\n", err)
		} else if count > 0 {
			fmt.Printf("[Scheduler] Cleaned up %d login logs older than %d days\n", count, logAutoDeleteDays)
		}
	}
}

func getConfigInt(svcCtx *svc.ServiceContext, key string, defaultVal int64) int64 {
	config, err := svcCtx.SystemConfigModel.FindByKey(context.Background(), key)
	if err != nil || config == nil {
		return defaultVal
	}
	val, err := strconv.ParseInt(config.ConfigValue, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
