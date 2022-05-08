// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	user "micro-todo/service/user/api/internal/handler/user"
	"micro-todo/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/info",
				Handler: user.UpdateInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/pwd",
				Handler: user.UpdatePwdHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/token",
				Handler: user.RefreshTokenHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Refresh.AccessSecret),
		rest.WithPrefix("/user"),
	)
}