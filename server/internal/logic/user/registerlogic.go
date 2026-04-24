package user

import (
	"context"
	"database/sql"
	"strings"

	"server/internal/model"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 检查注册开关
	config, err := l.svcCtx.SystemConfigModel.FindByKey(l.ctx, "allow_register")
	if err == nil && config.ConfigValue == "false" {
		return nil, xerr.NewCodeError(xerr.RegisterClosed)
	}

	// 2. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	// 3. 插入用户（直接 Insert，依赖 UNIQUE 约束避免 TOCTOU 竞态）
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     0,
		Status:   1,
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, xerr.NewCodeError(xerr.UserAlreadyExist)
		}
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	id, _ := result.LastInsertId()

	// 5. 记录操作日志
	_ , _ = l.svcCtx.OperationLogModel.Insert(l.ctx, &model.OperationLog{
		UserId:   sql.NullInt64{Int64: id, Valid: true},
		Username: req.Username,
		Module:   "user",
		Action:   "register",
		Method:   "POST",
		Status:   1,
	})

	return &types.RegisterResp{
		Id:       id,
		Username: req.Username,
	}, nil
}
