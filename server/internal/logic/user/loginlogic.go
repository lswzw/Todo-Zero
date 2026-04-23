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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 查找用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if err == model.ErrNotFound {
			_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
				Username: req.Username,
				Status:   0,
				Remark:   sql.NullString{String: "用户不存在", Valid: true},
			})
			return nil, xerr.NewCodeError(xerr.UserNotFoundError)
		}
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 2. 检查用户状态
	if user.Status == 0 {
		_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
			UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
			Username: req.Username,
			Status:   0,
			Remark:   sql.NullString{String: "用户已被禁用", Valid: true},
		})
		return nil, xerr.NewCodeError(xerr.UserDisabled)
	}

	// 3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
			UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
			Username: req.Username,
			Status:   0,
			Remark:   sql.NullString{String: "密码错误", Valid: true},
		})
		return nil, xerr.NewCodeError(xerr.PasswordError)
	}

	// 4. 生成 JWT Token
	now := time.Now()
	token, err := l.generateToken(user.Id, user.IsAdmin, now)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 5. 记录登录成功日志
	_ , _ = l.svcCtx.LoginLogModel.Insert(l.ctx, &model.LoginLog{
		UserId:   sql.NullInt64{Int64: user.Id, Valid: true},
		Username: req.Username,
		Status:   1,
		Remark:   sql.NullString{String: "登录成功", Valid: true},
	})

	return &types.LoginResp{
		Token:   token,
		IsAdmin: user.IsAdmin,
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
