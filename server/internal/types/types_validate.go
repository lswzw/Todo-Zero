package types

import (
	"fmt"
	"unicode/utf8"
)

// --- BatchTaskReq ---

func (r *BatchTaskReq) Validate() error {
	if len(r.Ids) == 0 {
		return fmt.Errorf("ids不能为空")
	}
	if len(r.Ids) > 100 {
		return fmt.Errorf("批量操作最多100条")
	}
	if r.Action == "" {
		return fmt.Errorf("action不能为空")
	}
	switch r.Action {
	case "complete", "undo", "delete":
	default:
		return fmt.Errorf("action必须是complete、undo或delete")
	}
	return nil
}

// --- ChangePasswordReq ---

func (r *ChangePasswordReq) Validate() error {
	if r.OldPassword == "" {
		return fmt.Errorf("原密码不能为空")
	}
	if err := validatePassword(r.NewPassword, "新密码"); err != nil {
		return err
	}
	return nil
}

// --- CreateCategoryReq ---

func (r *CreateCategoryReq) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("分类名称不能为空")
	}
	if utf8.RuneCountInString(r.Name) > 20 {
		return fmt.Errorf("分类名称最多20个字符")
	}
	if utf8.RuneCountInString(r.Color) > 20 {
		return fmt.Errorf("分类颜色最多20个字符")
	}
	return nil
}

// --- UpdateCategoryReq ---

func (r *UpdateCategoryReq) Validate() error {
	if r.Name != nil {
		if *r.Name == "" {
			return fmt.Errorf("分类名称不能为空")
		}
		if utf8.RuneCountInString(*r.Name) > 20 {
			return fmt.Errorf("分类名称最多20个字符")
		}
	}
	if r.Color != nil && utf8.RuneCountInString(*r.Color) > 20 {
		return fmt.Errorf("分类颜色最多20个字符")
	}
	if r.Icon != nil && utf8.RuneCountInString(*r.Icon) > 50 {
		return fmt.Errorf("分类图标最多50个字符")
	}
	if r.Sort != nil && *r.Sort < 0 {
		return fmt.Errorf("排序值无效")
	}
	return nil
}

// --- DeleteCategoryReq ---

func (r *DeleteCategoryReq) Validate() error {
	if r.Id <= 0 {
		return fmt.Errorf("分类ID无效")
	}
	return nil
}

// --- CreateTaskReq ---

func (r *CreateTaskReq) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("任务标题不能为空")
	}
	if utf8.RuneCountInString(r.Title) > 100 {
		return fmt.Errorf("任务标题最多100个字符")
	}
	if utf8.RuneCountInString(r.Content) > 1000 {
		return fmt.Errorf("任务内容最多1000个字符")
	}
	if r.Priority != 1 && r.Priority != 2 && r.Priority != 3 {
		return fmt.Errorf("优先级必须是1、2或3")
	}
	if r.CategoryId < 0 {
		return fmt.Errorf("分类ID无效")
	}
	return nil
}

// --- LoginLogReq ---

func (r *LoginLogReq) Validate() error {
	if r.Page < 1 {
		return fmt.Errorf("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return fmt.Errorf("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Username) > 20 {
		return fmt.Errorf("用户名最多20个字符")
	}
	return nil
}

// --- LoginReq ---

func (r *LoginReq) Validate() error {
	if err := validateUsername(r.Username, 1, 50); err != nil {
		return err
	}
	if r.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

// --- OperationLogReq ---

func (r *OperationLogReq) Validate() error {
	if r.Page < 1 {
		return fmt.Errorf("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return fmt.Errorf("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Action) > 20 {
		return fmt.Errorf("操作类型最多20个字符")
	}
	if utf8.RuneCountInString(r.Username) > 20 {
		return fmt.Errorf("用户名最多20个字符")
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
		return fmt.Errorf("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return fmt.Errorf("每页数量必须在1-100之间")
	}
	// Status: -1=全部, 0=待办, 2=已完成
	if r.Status != -1 && r.Status != 0 && r.Status != 2 {
		return fmt.Errorf("状态参数无效")
	}
	// Priority: -1=全部, 1=重要, 2=紧急, 3=普通
	if r.Priority != -1 && r.Priority != 1 && r.Priority != 2 && r.Priority != 3 {
		return fmt.Errorf("优先级参数无效")
	}
	if utf8.RuneCountInString(r.Keyword) > 50 {
		return fmt.Errorf("搜索关键词最多50个字符")
	}
	return nil
}

// --- UpdateConfigReq ---

func (r *UpdateConfigReq) Validate() error {
	if r.Key == "" {
		return fmt.Errorf("配置键不能为空")
	}
	if utf8.RuneCountInString(r.Key) > 50 {
		return fmt.Errorf("配置键最多50个字符")
	}
	if r.Value == "" {
		return fmt.Errorf("配置值不能为空")
	}
	if utf8.RuneCountInString(r.Value) > 500 {
		return fmt.Errorf("配置值最多500个字符")
	}
	return nil
}

// --- UpdateTaskReq ---

func (r *UpdateTaskReq) Validate() error {
	if r.Title != nil {
		if utf8.RuneCountInString(*r.Title) > 100 {
			return fmt.Errorf("任务标题最多100个字符")
		}
	}
	if r.Content != nil {
		if utf8.RuneCountInString(*r.Content) > 1000 {
			return fmt.Errorf("任务内容最多1000个字符")
		}
	}
	if r.Priority != nil {
		if *r.Priority != 1 && *r.Priority != 2 && *r.Priority != 3 {
			return fmt.Errorf("优先级必须是1、2或3")
		}
	}
	if r.CategoryId != nil && *r.CategoryId < 0 {
		return fmt.Errorf("分类ID无效")
	}
	return nil
}

// --- UserListReq ---

func (r *UserListReq) Validate() error {
	if r.Page < 1 {
		return fmt.Errorf("页码必须大于0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return fmt.Errorf("每页数量必须在1-100之间")
	}
	if utf8.RuneCountInString(r.Keyword) > 50 {
		return fmt.Errorf("搜索关键词最多50个字符")
	}
	return nil
}

// --- helpers ---

func validateUsername(username string, minLen, maxLen int) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	n := utf8.RuneCountInString(username)
	if n < minLen {
		return fmt.Errorf("用户名至少%d个字符", minLen)
	}
	if n > maxLen {
		return fmt.Errorf("用户名最多%d个字符", maxLen)
	}
	return nil
}

func validatePassword(password string, field string) error {
	if password == "" {
		return fmt.Errorf("%s不能为空", field)
	}
	n := utf8.RuneCountInString(password)
	if n < 6 {
		return fmt.Errorf("%s至少6个字符", field)
	}
	if n > 20 {
		return fmt.Errorf("%s最多20个字符", field)
	}
	// 复杂度要求：必须包含字母和数字
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
		return fmt.Errorf("%s必须同时包含字母和数字", field)
	}
	return nil
}
