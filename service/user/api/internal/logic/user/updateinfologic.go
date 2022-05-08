package user

import (
	"context"
	"encoding/json"
	"errors"

	"micro-todo/service/user/api/internal/svc"
	"micro-todo/service/user/api/internal/types"
	"micro-todo/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInfoLogic) UpdateInfo(req *types.UpdateInfoReq) (err error) {
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return errors.New("unexpected token data type")
	}

	if req.Username != "" {
		_, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
		switch err {
		case nil:
			return errors.New("existence user")
		case model.ErrNotFound:
		default:
			return err
		}
	}

	return l.svcCtx.UserModel.PartialUpdate(l.ctx, &model.User{
		Id:       uid,
		Username: req.Username,
		Gender:   req.Gender,
	})
}
