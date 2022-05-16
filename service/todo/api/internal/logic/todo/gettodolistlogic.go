package todo

import (
	"context"
	"fmt"

	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"
	"micro-todo/service/todo/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTodoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodoListLogic {
	return &GetTodoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTodoListLogic) GetTodoList(req *types.GetTodoListReq) (_ *types.GetTodoListReply, err error) {
	fmt.Println(req.KeyWorld)
	var todoList *[]model.Todo
	var count int

	if req.Type != "all" && req.KeyWorld == "" {
		todoList, count, err = l.svcCtx.TodoModel.FindListByDone(l.ctx, req.Type == "done", req.PageSize, req.PageNumber)
	}

	if req.Type == "all" && req.KeyWorld != "" {
		todoList, count, err = l.svcCtx.TodoModel.FindListByKeyword(l.ctx, req.KeyWorld, req.PageSize, req.PageNumber)
	}

	if req.Type != "all" && req.KeyWorld != "" {
		todoList, count, err = l.svcCtx.TodoModel.FindListByDoneAndKeyword(l.ctx, req.Type == "done", req.KeyWorld, req.PageSize, req.PageNumber)
	}

	if req.Type == "all" && req.KeyWorld == "" {
		todoList, count, err = l.svcCtx.TodoModel.FindList(l.ctx, req.PageSize, req.PageNumber)
	}

	if err != nil {
		return nil, err
	}
	var resp types.GetTodoListReply
	for _, todo := range *todoList {
		resp.TodoList = append(resp.TodoList, types.TodoReply{
			Id:      todo.Id,
			Title:   todo.Title,
			Content: todo.Content,
			Done:    todo.Done,
		})
	}
	resp.Count = count
	return &resp, nil
}
