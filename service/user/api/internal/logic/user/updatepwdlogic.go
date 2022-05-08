package user

import (
	"context"
	"encoding/json"
	"errors"

	"micro-todo/service/user/api/internal/svc"
	"micro-todo/service/user/api/internal/types"
	"micro-todo/service/user/model"
	"micro-todo/util/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePwdLogic {
	return &UpdatePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePwdLogic) UpdatePwd(req *types.UpdatePwdReq) (err error) {
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return errors.New("unexpected token data type")
	}

	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	switch err {
	case nil:
	case model.ErrNotFound:
		return errors.New("inexistence user")
	default:
		return err
	}

	if !bcrypt.Verify(userInfo.Password, req.OldPassword) {
		return errors.New("incorrect password")
	}

	hashed_pwd, err := bcrypt.Encrypt(req.NewPassword)
	if err != nil {
		return err
	}

	return l.svcCtx.UserModel.PartialUpdate(l.ctx, &model.User{
		Id:       uid,
		Password: hashed_pwd,
	})
}
