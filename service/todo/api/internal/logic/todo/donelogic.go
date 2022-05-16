package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoneLogic {
	return &DoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DoneLogic) Done(req *types.DoneReq) error {
	return l.svcCtx.TodoModel.SetTodoState(l.ctx, req.State, req.Id)
}
