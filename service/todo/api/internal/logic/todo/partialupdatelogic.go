package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"
	"micro-todo/service/todo/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PartialUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPartialUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PartialUpdateLogic {
	return &PartialUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PartialUpdateLogic) PartialUpdate(req *types.PartialUpdateReq) error {
	return l.svcCtx.TodoModel.PartialUpdate(l.ctx, &model.Todo{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
}
