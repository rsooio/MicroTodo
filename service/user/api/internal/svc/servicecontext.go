package svc

import (
	"micro-todo/service/user/api/internal/config"
	"micro-todo/service/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	userModel := model.NewUserModel(conn, c.CacheRedis)

	return &ServiceContext{
		Config:    c,
		UserModel: userModel,
	}
}
