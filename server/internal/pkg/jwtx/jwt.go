package jwtx

import (
	"encoding/json"

	"server/internal/pkg/xerr"
)

// GetUserIdFromCtx 从 JWT context 中获取用户ID
func GetUserIdFromCtx(ctx interface{ Value(any) any }) (int64, error) {
	val := ctx.Value("userId")
	if val == nil {
		return 0, xerr.NewCodeError(xerr.NoPermission)
	}
	switch v := val.(type) {
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, xerr.NewCodeError(xerr.NoPermission)
		}
		return n, nil
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, xerr.NewCodeError(xerr.NoPermission)
	}
}

// GetIsAdminFromCtx 从 JWT context 中获取是否管理员
func GetIsAdminFromCtx(ctx interface{ Value(any) any }) (int64, error) {
	val := ctx.Value("isAdmin")
	if val == nil {
		return 0, xerr.NewCodeError(xerr.AdminRequired)
	}
	switch v := val.(type) {
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, xerr.NewCodeError(xerr.AdminRequired)
		}
		return n, nil
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, xerr.NewCodeError(xerr.AdminRequired)
	}
}
