package todo

import (
	"net/http"

	"micro-todo/service/todo/api/internal/logic/todo"
	"micro-todo/service/todo/api/internal/svc"
	"micro-todo/service/todo/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTodoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTodoListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := todo.NewGetTodoListLogic(r.Context(), svcCtx)
		resp, err := l.GetTodoList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
