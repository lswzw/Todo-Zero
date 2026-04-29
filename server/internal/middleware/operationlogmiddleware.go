package middleware

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"server/internal/model"
	"server/internal/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

// OperationLogMiddleware automatically logs write operations (POST, PUT, PATCH, DELETE)
type OperationLogMiddleware struct {
	UserModel         model.UserModel
	OperationLogModel model.OperationLogModel
}

func NewOperationLogMiddleware(userModel model.UserModel, opLogModel model.OperationLogModel) *OperationLogMiddleware {
	return &OperationLogMiddleware{
		UserModel:         userModel,
		OperationLogModel: opLogModel,
	}
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only log write operations (POST, PUT, PATCH, DELETE)
		method := r.Method
		if method != http.MethodPost && method != http.MethodPut &&
			method != http.MethodPatch && method != http.MethodDelete {
			next(w, r)
			return
		}

		start := time.Now()

		// Execute the handler
		next(w, r)

		// Log the operation asynchronously (don't block the response)
		go func() {
			duration := time.Since(start).Milliseconds()

			userId := int64(0)
			username := "anonymous"
			if uid, err := jwtx.GetUserIdFromCtx(r.Context()); err == nil {
				userId = uid
			}
			// Look up username from DB
			if userId > 0 {
				if user, err := m.UserModel.FindOne(r.Context(), userId); err == nil {
					username = user.Username
				}
			}

			// Extract module and action from path
			module, action := parseRoute(r.URL.Path, method)

			_, _ = m.OperationLogModel.Insert(r.Context(), &model.OperationLog{
				UserId:   sql.NullInt64{Int64: userId, Valid: userId > 0},
				Username: username,
				Module:   module,
				Action:   action,
				Method:   method,
				Ip:       maskIP(r.RemoteAddr),
				Status:   1,
				Duration: duration,
			})

			logx.Infof("[OperationLog] %s %s by user:%s module:%s action:%s duration:%dms",
				method, r.URL.Path, username, module, action, duration)
		}()
	}
}

// maskIP masks the last octet of an IP address for privacy protection.
func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	// Handle IPv4 addresses
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) >= 4 {
			return parts[0] + "." + parts[1] + "." + parts[2] + ".x"
		}
	}
	// Handle IPv6 addresses
	if strings.Contains(ip, ":") {
		return "::1"
	}
	return "***.***.***.x"
}

// parseRoute extracts module and action from the URL path and method.
func parseRoute(path string, method string) (module string, action string) {
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/"), "/")

	// Determine module
	if len(parts) >= 1 {
		module = parts[0]
	}

	// Determine action based on method
	switch method {
	case http.MethodPost:
		action = "create"
	case http.MethodDelete:
		action = "delete"
	case http.MethodPut, http.MethodPatch:
		action = "update"
	default:
		action = strings.ToLower(method)
	}

	// Refine action based on path segments
	if len(parts) >= 2 {
		subPath := parts[len(parts)-1]
		switch {
		case strings.Contains(subPath, "toggle"):
			action = "toggle"
		case strings.Contains(subPath, "password"):
			action = "reset_password"
		case strings.Contains(subPath, "batch"):
			action = "batch"
		}
	}

	// Prefix action with sub-module for admin routes
	if module == "admin" && len(parts) >= 2 {
		action = parts[1] + "_" + action
	}

	return module, action
}
