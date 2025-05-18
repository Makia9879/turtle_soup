package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"turtle-soup/internal/config"
	"turtle-soup/internal/custom"
	"turtle-soup/internal/middleware"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
)

type ServiceContext struct {
	Config configurator.Configurator[config.Config]
	middleware.Middleware
	Custom  *custom.Custom
	SqlConn sqlx.SqlConn
}

func NewServiceContext(cc configurator.Configurator[config.Config]) *ServiceContext {
	sc := &ServiceContext{
		Config:     cc,
		Custom:     custom.New(),
		Middleware: middleware.New(),
		// TODO SqlConn
	}
	sc.SetConfigListener()
	return sc
}
