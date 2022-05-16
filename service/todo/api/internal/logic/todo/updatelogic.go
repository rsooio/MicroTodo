package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"
	"micro-todo/service/todo/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateReq) error {
	return l.svcCtx.TodoModel.ForceUpdate(l.ctx, &model.Todo{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
}
