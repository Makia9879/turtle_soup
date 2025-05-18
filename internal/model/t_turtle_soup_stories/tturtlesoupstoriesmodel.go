package t_turtle_soup_stories

import (
	"github.com/eddieowens/opts"
	"github.com/jzero-io/jzero-contrib/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TTurtleSoupStoriesModel = (*customTTurtleSoupStoriesModel)(nil)

type (
	// TTurtleSoupStoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTTurtleSoupStoriesModel.
	TTurtleSoupStoriesModel interface {
		tTurtleSoupStoriesModel
	}

	customTTurtleSoupStoriesModel struct {
		*defaultTTurtleSoupStoriesModel
	}
)

// NewTTurtleSoupStoriesModel returns a model for the database table.
func NewTTurtleSoupStoriesModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) TTurtleSoupStoriesModel {
	return &customTTurtleSoupStoriesModel{
		defaultTTurtleSoupStoriesModel: newTTurtleSoupStoriesModel(conn, op...),
	}
}
