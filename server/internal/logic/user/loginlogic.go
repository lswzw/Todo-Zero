package user

import (
	"context"
	"database/sql"
	"time"

	"server/internal/model"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	maxLoginAttempts     = 5
	lockDurationMinutes  = 15
)

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 查找用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if err == model.ErrNotFound {
			_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
				Username: req.Username,
				Status:   0,
				Remark:   "用户不存在",
			})
			return nil, xerr.NewCodeError(xerr.UserOrPasswordError)
		}
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 2. 检查账户是否被锁定
	if user.LockedUntil.Valid && user.LockedUntil.Time.After(time.Now()) {
		_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
			UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
			Username: req.Username,
			Status:   0,
			Remark:   "账户已被锁定",
		})
		return nil, xerr.NewCodeError(xerr.AccountLocked)
	}

	// 3. 检查用户状态
	if user.Status == 0 {
		_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
			UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
			Username: req.Username,
			Status:   0,
			Remark:   "用户已被禁用",
		})
		return nil, xerr.NewCodeError(xerr.UserOrPasswordError)
	}

	// 4. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
			UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
			Username: req.Username,
			Status:   0,
			Remark:   "密码错误",
		})
		
		// 记录失败次数并检查是否需要锁定
		_ = l.svcCtx.UserModel.IncrementFailedAttempts(l.ctx, user.Id, maxLoginAttempts, lockDurationMinutes)
		
		return nil, xerr.NewCodeError(xerr.UserOrPasswordError)
	}

	// 4. 生成 JWT Token
	now := time.Now()
	token, err := l.generateToken(user.Id, user.Role, now)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 5. 重置失败次数
	_ = l.svcCtx.UserModel.ResetFailedAttempts(l.ctx, user.Id)

	// 6. 记录登录成功日志
	_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
		UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
		Username: req.Username,
		Status:   1,
		Remark:   "登录成功",
	})

	return &types.LoginResp{
		Token:   token,
		IsAdmin: user.Role,
	}, nil
}

func (l *LoginLogic) generateToken(userId, isAdmin int64, now time.Time) (string, error) {
	claims := make(jwt.MapClaims)
	claims["userId"] = userId
	claims["isAdmin"] = isAdmin
	claims["exp"] = now.Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second).Unix()
	claims["iat"] = now.Unix()

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}
