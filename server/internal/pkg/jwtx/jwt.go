package jwtx

import (
	"encoding/json"

	"server/internal/pkg/xerr"
)

// GetUserIdFromCtx 从 JWT context 中获取用户ID
// 注意：此函数假设调用者已验证 JWT token 的有效性和过期时间
// JWT 验证由 go-zero 的 rest.WithJwt() 中间件在路由层完成
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
// 注意：此函数假设调用者已验证 JWT token 的有效性和过期时间
// JWT 验证由 go-zero 的 rest.WithJwt() 中间件在路由层完成
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
