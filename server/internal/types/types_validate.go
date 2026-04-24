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
	return nil
}

// --- LoginReq ---

func (r *LoginReq) Validate() error {
	if err := validateUsername(r.Username, 1, 100); err != nil {
		return err
	}
	if r.Password == "" {
		return fmt.Errorf("密码不能为空")
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

// --- UpdateConfigReq ---

func (r *UpdateConfigReq) Validate() error {
	if r.Key == "" {
		return fmt.Errorf("配置键不能为空")
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
	return nil
}
