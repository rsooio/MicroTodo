package middleware

import (
	"context"
	"micro-todo/service/user/rpc/types/user"
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}

func (m *AuthMiddleware) HandleWithUserRpc(userRpc user.UserClient) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
			if err != nil {
				httpx.Error(w, err)
				return
			}

			ctx := r.Context()
			reply, err := userRpc.Auth(ctx, &user.AuthReq{
				Token: token,
			})
			if err != nil {
				httpx.Error(w, err)
				return
			}

			ctx = context.WithValue(ctx, "uid", reply.Uid)

			next(w, r.WithContext(ctx))
		}
	}
}
