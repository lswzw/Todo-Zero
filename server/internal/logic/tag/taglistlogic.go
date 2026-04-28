package tag

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagListLogic) TagList(req *types.TagListReq) (resp *types.TagListResp, err error) {
	userId, err := jwtx.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.NoPermission)
	}
	tags, err := l.svcCtx.TagModel.FindList(l.ctx, userId, req.Keyword)
	if err != nil {
		return nil, err
	}

	resp = &types.TagListResp{
		List: make([]types.TagItem, 0, len(tags)),
	}
	for _, tag := range tags {
		resp.List = append(resp.List, types.TagItem{
			Id:       tag.Id,
			Name:     tag.Name,
			Color:    tag.Color,
			IsSystem: tag.IsSystem,
		})
	}

	return resp, nil
}
