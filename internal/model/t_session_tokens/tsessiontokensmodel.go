package t_session_tokens

import (
	"github.com/eddieowens/opts"
	"github.com/jzero-io/jzero-contrib/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TSessionTokensModel = (*customTSessionTokensModel)(nil)

type (
	// TSessionTokensModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTSessionTokensModel.
	TSessionTokensModel interface {
		tSessionTokensModel
	}

	customTSessionTokensModel struct {
		*defaultTSessionTokensModel
	}
)

// NewTSessionTokensModel returns a model for the database table.
func NewTSessionTokensModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) TSessionTokensModel {
	return &customTSessionTokensModel{
		defaultTSessionTokensModel: newTSessionTokensModel(conn, op...),
	}
}
