package logic

import (
	"context"

	"micro-todo/service/user/rpc/internal/svc"
	"micro-todo/service/user/rpc/types/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthLogic {
	return &AuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthLogic) Auth(in *user.AuthReq) (*user.AuthReply, error) {
	success := false
	msg := ""

	token, err := jwt.Parse(in.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})

	if err != nil {
		msg = err.Error()
	} else if !token.Valid {
		msg = "invalid auth token"
	} else {
		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			msg = "no auth params"
		} else {
			success = true
		}
	}

	return &user.AuthReply{
		Success: success,
		Msg:     msg,
	}, nil
}
