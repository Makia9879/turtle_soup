package t_activity_tokens

import (
	"github.com/eddieowens/opts"
	"github.com/jzero-io/jzero-contrib/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TActivityTokensModel = (*customTActivityTokensModel)(nil)

type (
	// TActivityTokensModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTActivityTokensModel.
	TActivityTokensModel interface {
		tActivityTokensModel
	}

	customTActivityTokensModel struct {
		*defaultTActivityTokensModel
	}
)

// NewTActivityTokensModel returns a model for the database table.
func NewTActivityTokensModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) TActivityTokensModel {
	return &customTActivityTokensModel{
		defaultTActivityTokensModel: newTActivityTokensModel(conn, op...),
	}
}
