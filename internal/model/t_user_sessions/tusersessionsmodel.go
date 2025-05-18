package t_user_sessions

import (
	"github.com/eddieowens/opts"
	"github.com/jzero-io/jzero-contrib/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TUserSessionsModel = (*customTUserSessionsModel)(nil)

type (
	// TUserSessionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTUserSessionsModel.
	TUserSessionsModel interface {
		tUserSessionsModel
	}

	customTUserSessionsModel struct {
		*defaultTUserSessionsModel
	}
)

// NewTUserSessionsModel returns a model for the database table.
func NewTUserSessionsModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) TUserSessionsModel {
	return &customTUserSessionsModel{
		defaultTUserSessionsModel: newTUserSessionsModel(conn, op...),
	}
}
