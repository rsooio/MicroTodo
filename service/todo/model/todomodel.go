package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TodoModel = (*customTodoModel)(nil)

type (
	// TodoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTodoModel.
	TodoModel interface {
		todoModel
		SecureUpdate(ctx context.Context, data *Todo) error
		ForceUpdate(ctx context.Context, data *Todo) error
		PartialUpdate(ctx context.Context, data *Todo) (err error)
		FindList(ctx context.Context, pageSize int, pageNum int) (_ *[]Todo, count int, err error)
		FindListByDone(ctx context.Context, done bool, pageSize int, pageNum int) (todoList *[]Todo, count int, err error)
		FindListByKeyword(ctx context.Context, keyword string, pageSize int, pageNum int) (todoList *[]Todo, count int, err error)
		FindListByDoneAndKeyword(ctx context.Context, done bool, keyWord string, pageSize int, pageNum int) (todoList *[]Todo, count int, err error)
		SetTodoState(ctx context.Context, state bool, id int64) error
		BatchSetTodoState(ctx context.Context, state bool) error
		SecureDelete(ctx context.Context, id int64) error
		BatchDelete(ctx context.Context, state bool) (sql.Result, error)
		AllDelete(ctx context.Context) (sql.Result, error)
	}

	customTodoModel struct {
		*defaultTodoModel
	}
)

// NewTodoModel returns a model for the database table.
func NewTodoModel(conn sqlx.SqlConn, c cache.CacheConf) TodoModel {
	return &customTodoModel{
		defaultTodoModel: newTodoModel(conn, c),
	}
}

func (m *customTodoModel) SecureUpdate(ctx context.Context, data *Todo) error {
	uid := ctx.Value("uid").(int64)

	todoIdKey := fmt.Sprintf("%s%v", cacheTodoIdPrefix, data.Id)
	todoTitleKey := fmt.Sprintf("%s%v", cacheTodoTitlePrefix, data.Title)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and `user_id` = ?", m.table, todoRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Title, data.Content, data.Done, data.UserId, data.Id, uid)
	}, todoIdKey, todoTitleKey)
	return err
}

func (m *customTodoModel) ForceUpdate(ctx context.Context, data *Todo) error {
	uid := ctx.Value("uid").(int64)

	todoInfo, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	} else if todoInfo.UserId != uid {
		return errors.New("access denied")
	}

	todoInfo.Title = data.Title
	todoInfo.Content = data.Content

	return m.Update(ctx, todoInfo)
}

func (m *customTodoModel) PartialUpdate(ctx context.Context, data *Todo) error {
	uid := ctx.Value("uid").(int64)

	todoInfo, err := m.FindOne(ctx, data.Id)
	if err != nil {
		return err
	} else if todoInfo.UserId != uid {
		return errors.New("access denied")
	}

	if data.Title != "" {
		todoInfo.Title = data.Title
	}
	if data.Content != "" {
		todoInfo.Content = data.Content
	}

	return m.Update(ctx, todoInfo)
}

func (m *customTodoModel) FindList(ctx context.Context, pageSize int, pageNum int) (_ *[]Todo, count int, err error) {
	uid := ctx.Value("uid").(int64)

	countQuery := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	err = m.QueryRowNoCacheCtx(ctx, &count, countQuery, uid)
	if err != nil {
		return nil, 0, err
	}

	var todoList []Todo
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit ? offset ?", todoRows, m.table)
	err = m.QueryRowsNoCacheCtx(ctx, &todoList, query, uid, pageSize, (pageNum-1)*pageSize)
	switch err {
	case nil:
		return &todoList, count, nil
	case sqlc.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customTodoModel) FindListByDone(ctx context.Context, done bool, pageSize int, pageNum int) (_ *[]Todo, count int, err error) {
	uid := ctx.Value("uid").(int64)

	countQuery := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `done` = ?", m.table)
	err = m.QueryRowNoCacheCtx(ctx, &count, countQuery, uid, 1)
	if err != nil {
		return nil, 0, err
	}

	var todoList []Todo
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `done` = ? limit ? offset ?", todoRows, m.table)
	err = m.QueryRowsNoCacheCtx(ctx, &todoList, query, uid, done, pageSize, (pageNum-1)*pageSize)
	switch err {
	case nil:
		return &todoList, count, nil
	case sqlc.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customTodoModel) FindListByKeyword(ctx context.Context, keyword string, pageSize int, pageNum int) (_ *[]Todo, count int, err error) {
	uid := ctx.Value("uid").(int64)

	countQuery := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `title` like ?", m.table)
	err = m.QueryRowNoCacheCtx(ctx, &count, countQuery, uid, "%"+keyword+"%")
	if err != nil {
		return nil, 0, err
	}

	var todoList []Todo
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `title` like ? limit ? offset ?", todoRows, m.table)
	err = m.QueryRowsNoCacheCtx(ctx, &todoList, query, uid, "%"+keyword+"%", pageSize, (pageNum-1)*pageSize)
	switch err {
	case nil:
		return &todoList, count, nil
	case sqlc.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customTodoModel) FindListByDoneAndKeyword(ctx context.Context, done bool, keyword string, pageSize int, pageNum int) (_ *[]Todo, count int, err error) {
	uid := ctx.Value("uid").(int64)

	countQuery := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `done` = ? and `title` like ?", m.table)
	err = m.QueryRowNoCacheCtx(ctx, &count, countQuery, uid, done, "%"+keyword+"%")
	if err != nil {
		return nil, 0, err
	}

	var todoList []Todo
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `done` = ? and `title` like ? limit ? offset ?", todoRows, m.table)
	err = m.QueryRowsNoCacheCtx(ctx, &todoList, query, uid, done, "%"+keyword+"%", pageSize, (pageNum-1)*pageSize)
	switch err {
	case nil:
		return &todoList, count, nil
	case sqlc.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customTodoModel) SetTodoState(ctx context.Context, state bool, id int64) error {
	uid := ctx.Value("uid").(int64)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set `done` = ? where `id` = ? and `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, state, id, uid)
	})
	return err
}

func (m *customTodoModel) BatchSetTodoState(ctx context.Context, state bool) error {
	uid := ctx.Value("uid").(int64)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set `done` = ? where `user_id` = ? and `done` = ?", m.table)
		return conn.ExecCtx(ctx, query, state, uid, !state)
	})
	return err
}

func (m *customTodoModel) SecureDelete(ctx context.Context, id int64) error {
	uid := ctx.Value("uid").(int64)

	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	} else if data.UserId != uid {
		return errors.New("access denied")
	}

	todoIdKey := fmt.Sprintf("%s%v", cacheTodoIdPrefix, id)
	todoTitleKey := fmt.Sprintf("%s%v", cacheTodoTitlePrefix, data.Title)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, todoIdKey, todoTitleKey)
	return err
}

func (m *customTodoModel) BatchDelete(ctx context.Context, state bool) (sql.Result, error) {
	uid := ctx.Value("uid").(int64)

	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `done` = ? and `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, state, uid)
	})
}

func (m *customTodoModel) AllDelete(ctx context.Context) (sql.Result, error) {
	uid := ctx.Value("uid").(int64)

	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, uid)
	})
}
