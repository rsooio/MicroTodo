package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLogic) BatchDelete(req *types.BatchDeleteReq) (err error) {
	if req.Type == "all" {
		_, err = l.svcCtx.TodoModel.AllDelete(l.ctx)
	} else {
		_, err = l.svcCtx.TodoModel.BatchDelete(l.ctx, req.Type == "done")
	}

	return err
}
