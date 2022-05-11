package user

import (
	"context"
	"errors"
	"time"

	"micro-todo/service/user/api/internal/svc"
	"micro-todo/service/user/api/internal/types"
	"micro-todo/service/user/model"
	"micro-todo/utils/bcrypt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginReply, err error) {
	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("inexistence user")
	default:
		return nil, err
	}

	if !bcrypt.Verify(userInfo.Password, req.Password) {
		return nil, errors.New("incorrect user or password")
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}
	refreshExpire := l.svcCtx.Config.Refresh.AccessExpire
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Refresh.AccessSecret, now, refreshExpire, userInfo.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginReply{
		Id:            userInfo.Id,
		Username:      userInfo.Username,
		Gender:        userInfo.Gender,
		AccessToken:   accessToken,
		AccessExpire:  now + accessExpire,
		RefreshToken:  refreshToken,
		RefreshExpire: now + refreshExpire,
	}, nil
}
