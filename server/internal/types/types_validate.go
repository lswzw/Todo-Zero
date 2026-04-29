package types

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"server/internal/pkg/xerr"
)

// --- SortTaskReq ---

func (r *SortTaskReq) Validate() error {
	if len(r.Orders) == 0 {
		return xerr.NewValidationError("排序数据不能为空")
	}
	if len(r.Orders) > 100 {
		return xerr.NewValidationError("排序数据最多100条")
	}
	for i, item := range r.Orders {
		if item.Id <= 0 {
			return xerr.NewValidationError(fmt.Sprintf("第%d项ID无效", i+1))
		}
		if item.SortOrder < 0 {
			return xerr.NewValidationError(fmt.Sprintf("第%d项排序值无效", i+1))
		}
	}
	return nil
}

// --- BatchTaskReq ---

func (r *BatchTaskReq) Validate() error {
	if len(r.Ids) == 0 {
		return xerr.NewValidationError("ids不能为空")
	}
	if len(r.Ids) > 100 {
		return xerr.NewValidationError("批量操作最多100条")
	}
	if r.Action == "" {
		return xerr.NewValidationError("action不能为空")
	}
	switch r.Action {
	case "complete", "undo", "delete", "restore":
	default:
		return xerr.NewValidationError("action必须是complete、undo、delete或restore")
	}
	return nil
}

// --- ChangePasswordReq ---

func (r *ChangePasswordReq) Validate() error {
	if r.OldPassword == "" {
		return xerr.NewValidationError("原密码不能为空")
	}
	if err := validatePassword(r.NewPassword, "新密码"); err != nil {
		return err
	}
	return nil
}

// --- CreateCategoryReq ---

func (r *CreateCategoryReq) Validate() error {
	if r.Name == "" {
		return xerr.NewValidationError("分类名称不能为空")
	}
	if utf8.RuneCountInString(r.Name) > 20 {
		return xerr.NewValidationError("分类名称最多20个字符")
	}
	if utf8.RuneCountInString(r.Color) > 20 {
		return xerr.NewValidationError("分类颜色最多20个字符")
	}
	return nil
}

// --- UpdateCategoryReq ---

func (r *UpdateCategoryReq) Validate() error {
	if r.Name != nil {
		if *r.Name == "" {
			return xerr.NewValidationError("分类名称不能为空")
		}
		if utf8.RuneCountInString(*r.Name) > 20 {
			return xerr.NewValidationError("分类名称最多20个字符")
		}
	}
	if r.Color != nil && utf8.RuneCountInString(*r.Color) > 20 {
		return xerr.NewValidationError("分类颜色最多20个字符")
	}
	if r.Icon != nil && utf8.RuneCountInString(*r.Icon) > 50 {
		return xerr.NewValidationError("分类图标最多50个字符")
	}
	if r.Sort != nil && *r.Sort < 0 {
		return xerr.NewValidationError("排序值无效")
	}
	return nil
}

// --- DeleteCategoryReq ---

func (r *DeleteCategoryReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("分类ID无效")
	}
	return nil
}

// --- CreateTaskReq ---

func (r *CreateTaskReq) Validate() error {
	if r.Title == "" {
		return xerr.NewValidationError("任务标题不能为空")
	}
	if utf8.RuneCountInString(r.Title) > 100 {
		return xerr.NewValidationError("任务标题最多100个字符")
	}
	if utf8.RuneCountInString(r.Content) > 1000 {
		return xerr.NewValidationError("任务内容最多1000个字符")
	}
	if r.Priority != 1 && r.Priority != 2 && r.Priority != 3 {
		return xerr.NewValidationError("优先级必须是1、2或3")
	}
	if r.CategoryId < 0 {
		return xerr.NewValidationError("分类ID无效")
	}
	if utf8.RuneCountInString(r.Tags) > 200 {
		return xerr.NewValidationError("标签最多200个字符")
	}
	if r.StartTime != "" {
		if _, err := time.Parse("2006-01-02 15:04", r.StartTime); err != nil {
			return xerr.NewValidationError("开始时间格式无效")
		}
	}
	if r.EndTime != "" {
		if _, err := time.Parse("2006-01-02 15:04", r.EndTime); err != nil {
			return xerr.NewValidationError("截止时间格式无效")
		}
	}
	if r.Reminder != "" {
		if _, err := time.Parse("2006-01-02 15:04", r.Reminder); err != nil {
			return xerr.NewValidationError("提醒时间格式无效")
		}
	}
	return nil
}

// --- LoginLogReq ---

func (r *LoginLogReq) Validate() error {
	if r.Page < 1 {
		return xerr.NewValidationError("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return xerr.NewValidationError("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Username) > 20 {
		return xerr.NewValidationError("用户名最多20个字符")
	}
	return nil
}

// --- LoginReq ---

func (r *LoginReq) Validate() error {
	if err := validateUsername(r.Username, 1, 50); err != nil {
		return err
	}
	if r.Password == "" {
		return xerr.NewValidationError("密码不能为空")
	}
	return nil
}

// --- OperationLogReq ---

func (r *OperationLogReq) Validate() error {
	if r.Page < 1 {
		return xerr.NewValidationError("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return xerr.NewValidationError("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Action) > 20 {
		return xerr.NewValidationError("操作类型最多20个字符")
	}
	if utf8.RuneCountInString(r.Username) > 20 {
		return xerr.NewValidationError("用户名最多20个字符")
	}
	return nil
}

// --- RegisterReq ---

func (r *RegisterReq) Validate() error {
	if err := validateUsername(r.Username, 3, 20); err != nil {
		return err
	}
	if err := validatePassword(r.Password, "密码"); err != nil {
		return err
	}
	return nil
}

// --- ResetPasswordReq ---

func (r *ResetPasswordReq) Validate() error {
	if err := validatePassword(r.NewPassword, "新密码"); err != nil {
		return err
	}
	return nil
}

// --- TaskListReq ---

func (r *TaskListReq) Validate() error {
	if r.Page < 1 {
		return xerr.NewValidationError("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return xerr.NewValidationError("每页数量必须在1-100之间")
	}
	// Status: -1=全部, 0=待办, 2=已完成
	if r.Status != -1 && r.Status != 0 && r.Status != 2 {
		return xerr.NewValidationError("状态参数无效")
	}
	// Priority: -1=全部, 1=重要, 2=紧急, 3=普通
	if r.Priority != -1 && r.Priority != 1 && r.Priority != 2 && r.Priority != 3 {
		return xerr.NewValidationError("优先级参数无效")
	}
	if utf8.RuneCountInString(r.Keyword) > 50 {
		return xerr.NewValidationError("搜索关键词最多50个字符")
	}
	return nil
}

// --- UpdateConfigReq ---

func (r *UpdateConfigReq) Validate() error {
	if r.Key == "" {
		return xerr.NewValidationError("配置键不能为空")
	}
	if utf8.RuneCountInString(r.Key) > 50 {
		return xerr.NewValidationError("配置键最多50个字符")
	}
	if r.Value == "" {
		return xerr.NewValidationError("配置值不能为空")
	}
	if utf8.RuneCountInString(r.Value) > 500 {
		return xerr.NewValidationError("配置值最多500个字符")
	}
	return nil
}

// --- UpdateTaskReq ---

func (r *UpdateTaskReq) Validate() error {
	if r.Title != nil {
		if utf8.RuneCountInString(*r.Title) > 100 {
			return xerr.NewValidationError("任务标题最多100个字符")
		}
	}
	if r.Content != nil {
		if utf8.RuneCountInString(*r.Content) > 1000 {
			return xerr.NewValidationError("任务内容最多1000个字符")
		}
	}
	if r.Priority != nil {
		if *r.Priority != 1 && *r.Priority != 2 && *r.Priority != 3 {
			return xerr.NewValidationError("优先级必须是1、2或3")
		}
	}
	if r.CategoryId != nil && *r.CategoryId < 0 {
		return xerr.NewValidationError("分类ID无效")
	}
	if r.Tags != nil && utf8.RuneCountInString(*r.Tags) > 200 {
		return xerr.NewValidationError("标签最多200个字符")
	}
	if r.StartTime != nil && *r.StartTime != "" {
		if _, err := time.Parse("2006-01-02 15:04", *r.StartTime); err != nil {
			return xerr.NewValidationError("开始时间格式无效")
		}
	}
	if r.EndTime != nil && *r.EndTime != "" {
		if _, err := time.Parse("2006-01-02 15:04", *r.EndTime); err != nil {
			return xerr.NewValidationError("截止时间格式无效")
		}
	}
	if r.Reminder != nil && *r.Reminder != "" {
		if _, err := time.Parse("2006-01-02 15:04", *r.Reminder); err != nil {
			return xerr.NewValidationError("提醒时间格式无效")
		}
	}
	return nil
}

// --- UserListReq ---

func (r *UserListReq) Validate() error {
	if r.Page < 1 {
		return xerr.NewValidationError("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return xerr.NewValidationError("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Keyword) > 50 {
		return xerr.NewValidationError("搜索关键词最多50个字符")
	}
	return nil
}

// --- helpers ---

func validateUsername(username string, minLen, maxLen int) error {
	if username == "" {
		return xerr.NewValidationError("用户名不能为空")
	}
	n := utf8.RuneCountInString(username)
	if n < minLen {
		return xerr.NewValidationError(fmt.Sprintf("用户名至少%d个字符", minLen))
	}
	if n > maxLen {
		return xerr.NewValidationError(fmt.Sprintf("用户名最多%d个字符", maxLen))
	}
	return nil
}

func validatePassword(password string, field string) error {
	if password == "" {
		return xerr.NewValidationError(fmt.Sprintf("%s不能为空", field))
	}
	n := utf8.RuneCountInString(password)
	if n < 6 {
		return xerr.NewValidationError(fmt.Sprintf("%s至少6个字符", field))
	}
	if n > 20 {
		return xerr.NewValidationError(fmt.Sprintf("%s最多20个字符", field))
	}
	hasLetter := false
	hasDigit := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return xerr.NewValidationError(fmt.Sprintf("%s必须同时包含字母和数字", field))
	}
	return nil
}

// --- DownloadBackupReq ---

func (r *DownloadBackupReq) Validate() error {
	if err := validateFileName(r.FileName); err != nil {
		return err
	}
	return nil
}

// --- RestoreBackupReq ---

func (r *RestoreBackupReq) Validate() error {
	if err := validateFileName(r.FileName); err != nil {
		return err
	}
	return nil
}

// --- DeleteTaskReq ---

func (r *DeleteTaskReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("任务ID无效")
	}
	return nil
}

// --- RestoreTaskReq ---

func (r *RestoreTaskReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("任务ID无效")
	}
	return nil
}

// --- PermanentDeleteTaskReq ---

func (r *PermanentDeleteTaskReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("任务ID无效")
	}
	return nil
}

// --- DeleteUserReq ---

func (r *DeleteUserReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("用户ID无效")
	}
	return nil
}

// --- TaskDetailReq ---

func (r *TaskDetailReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("任务ID无效")
	}
	return nil
}

// --- ToggleTaskReq ---

func (r *ToggleTaskReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("任务ID无效")
	}
	return nil
}

// --- ToggleUserStatusReq ---

func (r *ToggleUserStatusReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("用户ID无效")
	}
	return nil
}

// --- CreateTagReq ---

func (r *CreateTagReq) Validate() error {
	if r.Name == "" {
		return xerr.NewValidationError("标签名称不能为空")
	}
	if utf8.RuneCountInString(r.Name) > 20 {
		return xerr.NewValidationError("标签名称最多20个字符")
	}
	if utf8.RuneCountInString(r.Color) > 20 {
		return xerr.NewValidationError("标签颜色最多20个字符")
	}
	return nil
}

// --- UpdateTagReq ---

func (r *UpdateTagReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("标签ID无效")
	}
	if r.Name != nil && utf8.RuneCountInString(*r.Name) > 20 {
		return xerr.NewValidationError("标签名称最多20个字符")
	}
	if r.Color != nil && utf8.RuneCountInString(*r.Color) > 20 {
		return xerr.NewValidationError("标签颜色最多20个字符")
	}
	return nil
}

// --- DeleteTagReq ---

func (r *DeleteTagReq) Validate() error {
	if r.Id <= 0 {
		return xerr.NewValidationError("标签ID无效")
	}
	return nil
}

// --- TagListReq ---

func (r *TagListReq) Validate() error {
	if utf8.RuneCountInString(r.Keyword) > 50 {
		return xerr.NewValidationError("搜索关键词最多50个字符")
	}
	return nil
}

// validateFileName validates backup file name for security
func validateFileName(fileName string) error {
	if fileName == "" {
		return xerr.NewValidationError("文件名不能为空")
	}
	if len(fileName) > 255 {
		return xerr.NewValidationError("文件名过长")
	}
	if strings.Contains(fileName, "..") {
		return xerr.NewValidationError("文件名包含非法字符")
	}
	if strings.HasPrefix(fileName, "/") || strings.HasPrefix(fileName, "\\") {
		return xerr.NewValidationError("文件名不能是绝对路径")
	}
	if !strings.HasSuffix(fileName, ".bak") && !strings.HasSuffix(fileName, ".BAK") {
		return xerr.NewValidationError("必须是.bak文件")
	}
	return nil
}
