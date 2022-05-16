package logic

import (
	"context"
	"errors"

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

func (l *AuthLogic) Auth(in *user.AuthReq) (out *user.AuthReply, err error) {
	token, err := jwt.ParseWithClaims(in.Token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.AuthInfo.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid auth token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no auth params")
	}

	uid, ok := claims["uid"]
	if !ok {
		return nil, errors.New("auth params error 1")
	}

	userId, ok := uid.(float64)
	if !ok {
		return nil, errors.New("auth params error 2")
	}

	return &user.AuthReply{Uid: int64(userId)}, nil
}
