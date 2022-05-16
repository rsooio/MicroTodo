package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDoneLogic {
	return &BatchDoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDoneLogic) BatchDone(req *types.BatchDoneReq) error {
	return l.svcCtx.TodoModel.BatchSetTodoState(l.ctx, req.State)
}
