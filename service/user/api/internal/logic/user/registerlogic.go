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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *RegisterLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.LoginReply, err error) {
	hashed_pwd, err := bcrypt.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	switch err {
	case nil:
		return nil, errors.New("existence user")
	case model.ErrNotFound:
	default:
		return nil, err
	}

	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: req.Username,
		Password: hashed_pwd,
		Gender:   req.Gender,
	})
	switch err {
	case nil:
	default:
		return nil, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, userId)
	if err != nil {
		return nil, err
	}
	refreshExpire := l.svcCtx.Config.Refresh.AccessExpire
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Refresh.AccessSecret, now, refreshExpire, userId)
	if err != nil {
		return nil, err
	}

	return &types.LoginReply{
		Id:            userId,
		Username:      req.Username,
		Gender:        req.Gender,
		AccessToken:   accessToken,
		AccessExpire:  now + accessExpire,
		RefreshToken:  refreshToken,
		RefreshExpire: now + refreshExpire,
	}, nil
}
