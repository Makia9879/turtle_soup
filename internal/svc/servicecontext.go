package svc

import (
	"fmt"
	"github.com/jzero-io/jzero-contrib/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"turtle-soup/internal/errors"
	"turtle-soup/internal/model"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"turtle-soup/internal/config"
	"turtle-soup/internal/custom"
	"turtle-soup/internal/middleware"
)

type ServiceContext struct {
	Config configurator.Configurator[config.Config]
	middleware.Middleware
	Custom *custom.Custom

	SqlConn   sqlx.SqlConn
	Model     model.Model
	Cache     cache.Cache
	RedisConn *redis.Redis
}

func NewServiceContext(cc configurator.Configurator[config.Config]) *ServiceContext {
	cfg, err := cc.GetConfig()
	logx.Must(err)

	sqlConn, err := initSqlConn(cfg)
	logx.Must(err)

	// 初始化 Cache
	redisConn, err := redis.NewRedis(cfg.RedisConf)
	logx.Must(err)
	cacheIns := cache.NewRedisNode(redisConn, errors.ErrCacheNotFound)

	sc := &ServiceContext{
		Config:     cc,
		Custom:     custom.New(),
		Middleware: middleware.New(),
		SqlConn:    sqlConn,
		Model:      model.NewModel(sqlConn),
		Cache:      cacheIns,
		RedisConn:  redisConn,
	}
	sc.SetConfigListener()
	return sc
}

func initSqlConn(c config.Config) (sqlx.SqlConn, error) {
	mysqlCfg := c.MySQLConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.IP,
		mysqlCfg.Port,
		mysqlCfg.DBName)
	return sqlx.NewMysql(dsn), nil
}
