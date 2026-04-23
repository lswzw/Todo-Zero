package jwtx

import "github.com/golang-jwt/jwt/v4"

// GetUserIdFromCtx 从 JWT context 中获取用户ID
func GetUserIdFromCtx(claims jwt.MapClaims) int64 {
	if userId, ok := claims["userId"]; ok {
		if id, ok := userId.(float64); ok {
			return int64(id)
		}
	}
	return 0
}

// GetIsAdminFromCtx 从 JWT context 中获取是否管理员
func GetIsAdminFromCtx(claims jwt.MapClaims) int64 {
	if isAdmin, ok := claims["isAdmin"]; ok {
		if v, ok := isAdmin.(float64); ok {
			return int64(v)
		}
	}
	return 0
}
