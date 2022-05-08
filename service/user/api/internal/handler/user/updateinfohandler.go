package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"micro-todo/service/user/api/internal/logic/user"
	"micro-todo/service/user/api/internal/svc"
	"micro-todo/service/user/api/internal/types"
)

func UpdateInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUpdateInfoLogic(r.Context(), svcCtx)
		err := l.UpdateInfo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
