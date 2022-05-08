package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"micro-todo/service/user/api/internal/svc"
	"micro-todo/service/user/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *RefreshTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken() (resp *types.TokenRefreshReply, err error) {
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return nil, errors.New("unexpected token data type")
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, uid)
	if err != nil {
		return nil, err
	}
	refreshExpire := l.svcCtx.Config.Refresh.AccessExpire
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Refresh.AccessSecret, now, refreshExpire, uid)
	if err != nil {
		return nil, err
	}

	return &types.TokenRefreshReply{
		AccessToken:   accessToken,
		AccessExpire:  now + accessExpire,
		RefreshToken:  refreshToken,
		RefreshExpire: now + refreshExpire,
	}, nil
}
