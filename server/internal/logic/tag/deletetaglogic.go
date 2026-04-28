package tag

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTagLogic {
	return &DeleteTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTagLogic) DeleteTag(req *types.DeleteTagReq) (resp *types.DeleteTagResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if tag.UserId != userId {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}

	if tag.IsSystem == 1 {
		return nil, xerr.NewCodeError(xerr.RequestParamError)
	}

	err = l.svcCtx.TagModel.Delete(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.DeleteTagResp{Success: true}, nil
}
