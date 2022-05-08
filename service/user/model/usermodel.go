package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		PartialUpdate(ctx context.Context, data *User) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c),
	}
}

func (m *customUserModel) PartialUpdate(ctx context.Context, data *User) error {
	userInfo, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	}

	if data.Password != "" {
		userInfo.Password = data.Password
	}
	if data.Username != "" {
		userInfo.Username = data.Username
	}
	if data.Gender != "" {
		userInfo.Username = data.Gender
	}

	return m.Update(ctx, userInfo)
}
