// Code generated by goctl. Templates Edited by jzero. DO NOT EDIT.

package t_user_sessions

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eddieowens/opts"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero-contrib/condition"
	"github.com/jzero-io/jzero-contrib/modelx"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	tUserSessionsFieldNames          = builder.RawFieldNames(&TUserSessions{})
	tUserSessionsRows                = strings.Join(tUserSessionsFieldNames, ",")
	tUserSessionsRowsExpectAutoSet   = strings.Join(stringx.Remove(tUserSessionsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	tUserSessionsRowsWithPlaceHolder = strings.Join(stringx.Remove(tUserSessionsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTUserSessionsIdPrefix = "cache:tUserSessions:id:"
)

type (
	tUserSessionsModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *TUserSessions) (sql.Result, error)
		InsertWithCache(ctx context.Context, session sqlx.Session, data *TUserSessions) (sql.Result, error)
		FindOne(ctx context.Context, session sqlx.Session, id uint64) (*TUserSessions, error)
		FindOneWithCache(ctx context.Context, session sqlx.Session, id uint64) (*TUserSessions, error)
		Update(ctx context.Context, session sqlx.Session, data *TUserSessions) error
		UpdateWithCache(ctx context.Context, session sqlx.Session, data *TUserSessions) error
		Delete(ctx context.Context, session sqlx.Session, id uint64) error
		DeleteWithCache(ctx context.Context, session sqlx.Session, id uint64) error

		// custom interface generated by jzero
		BulkInsert(ctx context.Context, session sqlx.Session, datas []*TUserSessions) error
		FindByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*TUserSessions, error)
		FindOneByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) (*TUserSessions, error)
		PageByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*TUserSessions, int64, error)
		UpdateFieldsByCondition(ctx context.Context, session sqlx.Session, field map[string]any, conds ...condition.Condition) error
		DeleteByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) error
	}

	defaultTUserSessionsModel struct {
		cachedConn sqlc.CachedConn
		conn       sqlx.SqlConn
		table      string
	}

	TUserSessions struct {
		Id           uint64       `db:"id"`
		SessionToken string       `db:"session_token"`
		StoryId      int64        `db:"story_id"`
		IsCompleted  int64        `db:"is_completed"`
		CompletedAt  sql.NullTime `db:"completed_at"`
		CreatedAt    time.Time    `db:"created_at"`
	}
)

func newTUserSessionsModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) *defaultTUserSessionsModel {
	o := opts.DefaultApply(op...)
	var cachedConn sqlc.CachedConn
	if len(o.CacheConf) > 0 {
		cachedConn = sqlc.NewConn(conn, o.CacheConf, o.CacheOpts...)
	}
	if o.CachedConn != nil {
		cachedConn = *o.CachedConn
	}
	return &defaultTUserSessionsModel{
		cachedConn: cachedConn,
		conn:       conn,
		table:      "`T_user_sessions`",
	}
}
func (m *defaultTUserSessionsModel) Delete(ctx context.Context, session sqlx.Session, id uint64) error {
	sb := sqlbuilder.DeleteFrom(m.table)
	sb.Where(sb.EQ("`id`", id))
	statement, args := sb.Build()
	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}

func (m *defaultTUserSessionsModel) DeleteWithCache(ctx context.Context, session sqlx.Session, id uint64) error {
	tUserSessionsIdKey := fmt.Sprintf("%s%v", cacheTUserSessionsIdPrefix, id)
	_, err := m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		sb := sqlbuilder.DeleteFrom(m.table)
		sb.Where(sb.EQ("`id`", id))
		statement, args := sb.Build()
		if session != nil {
			return session.ExecCtx(ctx, statement, args...)
		}
		return conn.ExecCtx(ctx, statement, args...)
	}, tUserSessionsIdKey)
	return err
}

func (m *defaultTUserSessionsModel) FindOne(ctx context.Context, session sqlx.Session, id uint64) (*TUserSessions, error) {
	sb := sqlbuilder.Select(tUserSessionsRows).From(m.table)
	sb.Where(sb.EQ("`id`", id))
	sb.Limit(1)
	sql, args := sb.Build()
	var resp TUserSessions
	var err error
	if session != nil {
		err = session.QueryRowCtx(ctx, &resp, sql, args...)
	} else {
		err = m.conn.QueryRowCtx(ctx, &resp, sql, args...)
	}
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUserSessionsModel) FindOneWithCache(ctx context.Context, session sqlx.Session, id uint64) (*TUserSessions, error) {
	tUserSessionsIdKey := fmt.Sprintf("%s%v", cacheTUserSessionsIdPrefix, id)
	var resp TUserSessions
	err := m.cachedConn.QueryRowCtx(ctx, &resp, tUserSessionsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		sb := sqlbuilder.Select(tUserSessionsRows).From(m.table)
		sb.Where(sb.EQ("`id`", id))
		sql, args := sb.Build()
		if session != nil {
			return session.QueryRowCtx(ctx, v, sql, args...)
		}
		return conn.QueryRowCtx(ctx, v, sql, args...)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUserSessionsModel) Insert(ctx context.Context, session sqlx.Session, data *TUserSessions) (sql.Result, error) {
	statement, args := sqlbuilder.NewInsertBuilder().
		InsertInto(m.table).
		Cols(tUserSessionsRowsExpectAutoSet).
		Values(data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt).Build()
	if session != nil {
		return session.ExecCtx(ctx, statement, args...)
	}
	return m.conn.ExecCtx(ctx, statement, args...)
}

func (m *defaultTUserSessionsModel) InsertWithCache(ctx context.Context, session sqlx.Session, data *TUserSessions) (sql.Result, error) {
	tUserSessionsIdKey := fmt.Sprintf("%s%v", cacheTUserSessionsIdPrefix, data.Id)
	statement, args := sqlbuilder.NewInsertBuilder().
		InsertInto(m.table).
		Cols(tUserSessionsRowsExpectAutoSet).
		Values(data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt).Build()
	return m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, statement, args...)
		}
		return conn.ExecCtx(ctx, statement, args...)
	}, tUserSessionsIdKey)
}
func (m *defaultTUserSessionsModel) Update(ctx context.Context, session sqlx.Session, data *TUserSessions) error {
	sb := sqlbuilder.Update(m.table)
	split := strings.Split(tUserSessionsRowsExpectAutoSet, ",")
	var assigns []string
	for _, s := range split {
		assigns = append(assigns, sb.Assign(s, nil))
	}
	sb.Set(assigns...)
	sb.Where(sb.EQ("`id`", nil))
	statement, _ := sb.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt, data.Id)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt, data.Id)
	}
	return err
}

func (m *defaultTUserSessionsModel) UpdateWithCache(ctx context.Context, session sqlx.Session, data *TUserSessions) error {
	tUserSessionsIdKey := fmt.Sprintf("%s%v", cacheTUserSessionsIdPrefix, data.Id)
	_, err := m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		sb := sqlbuilder.Update(m.table)
		split := strings.Split(tUserSessionsRowsExpectAutoSet, ",")
		var assigns []string
		for _, s := range split {
			assigns = append(assigns, sb.Assign(s, nil))
		}
		sb.Set(assigns...)
		sb.Where(sb.EQ("`id`", nil))
		statement, _ := sb.Build()
		if session != nil {
			return session.ExecCtx(ctx, statement, data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt, data.Id)
		}
		return conn.ExecCtx(ctx, statement, data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt, data.Id)
	}, tUserSessionsIdKey)
	return err
}

func (m *defaultTUserSessionsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTUserSessionsIdPrefix, primary)
}

func (m *defaultTUserSessionsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	sb := sqlbuilder.Select(tUserSessionsRows).From(m.table)
	sb.Where(sb.EQ("`id`", primary))
	sql, args := sb.Build()
	return conn.QueryRowCtx(ctx, v, sql, args...)
}

func (m *defaultTUserSessionsModel) tableName() string {
	return m.table
}

func (m *customTUserSessionsModel) BulkInsert(ctx context.Context, session sqlx.Session, datas []*TUserSessions) error {
	sb := sqlbuilder.InsertInto(m.table)
	sb.Cols(tUserSessionsRowsExpectAutoSet)
	for _, data := range datas {
		sb.Values(data.SessionToken, data.StoryId, data.IsCompleted, data.CompletedAt)
	}
	statement, args := sb.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}

func (m *customTUserSessionsModel) FindByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*TUserSessions, error) {
	sb := sqlbuilder.Select(tUserSessionsFieldNames...).From(m.table)
	builder := condition.Select(*sb, conds...)
	statement, args := builder.Build()

	var resp []*TUserSessions
	var err error

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customTUserSessionsModel) FindOneByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) (*TUserSessions, error) {
	sb := sqlbuilder.Select(tUserSessionsFieldNames...).From(m.table)

	builder := condition.Select(*sb, conds...)
	builder.Limit(1)
	statement, args := builder.Build()

	var resp TUserSessions
	var err error

	if session != nil {
		err = session.QueryRowCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customTUserSessionsModel) PageByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*TUserSessions, int64, error) {
	sb := sqlbuilder.Select(tUserSessionsFieldNames...).From(m.table)
	countsb := sqlbuilder.Select("count(*)").From(m.table)

	builder := condition.Select(*sb, conds...)

	var countConds []condition.Condition
	for _, cond := range conds {
		if cond.Operator != condition.Limit && cond.Operator != condition.Offset {
			countConds = append(countConds, cond)
		}
	}
	countBuilder := condition.Select(*countsb, countConds...)

	var resp []*TUserSessions
	var err error

	statement, args := builder.Build()

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, 0, err
	}

	var total int64
	statement, args = countBuilder.Build()
	if session != nil {
		err = session.QueryRowCtx(ctx, &total, statement, args...)
	} else {
		err = m.conn.QueryRowCtx(ctx, &total, statement, args...)
	}
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func (m *customTUserSessionsModel) UpdateFieldsByCondition(ctx context.Context, session sqlx.Session, field map[string]any, conds ...condition.Condition) error {
	if field == nil {
		return nil
	}

	sb := sqlbuilder.Update(m.table)
	builder := condition.Update(*sb, conds...)

	var assigns []string
	for key, value := range field {
		assigns = append(assigns, sb.Assign(key, value))
	}
	builder.Set(assigns...)

	statement, args := builder.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *customTUserSessionsModel) DeleteByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) error {
	if len(conds) == 0 {
		return nil
	}
	sb := sqlbuilder.DeleteFrom(m.table)
	builder := condition.Delete(*sb, conds...)
	statement, args := builder.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}
