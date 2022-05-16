package todo

import (
	"context"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"
	"micro-todo/service/todo/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateReq) (err error) {
	uid := l.ctx.Value("uid").(int64)

	_, err = l.svcCtx.TodoModel.Insert(l.ctx, &model.Todo{
		Title:   req.Title,
		Content: req.Content,
		UserId:  uid,
	})

	return
}
