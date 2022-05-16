package svc

import (
	"micro-todo/service/todo/api/internal/config"
	"micro-todo/service/todo/api/internal/middleware"
	"micro-todo/service/todo/model"
	"micro-todo/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Auth      rest.Middleware
	UserRpc   user.UserClient
	TodoModel model.TodoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRpc := user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn())
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	todoModel := model.NewTodoModel(conn, c.CacheRedis)

	return &ServiceContext{
		Config:    c,
		Auth:      middleware.NewAuthMiddleware().HandleWithUserRpc(userRpc),
		UserRpc:   userRpc,
		TodoModel: todoModel,
	}
}
