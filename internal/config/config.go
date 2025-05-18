package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	Rest RestConf
	Log  LogConf

	Banner BannerConf

	MySQLConf MySQLConf       `json:"mysql,optional"`
	RedisConf redis.RedisConf `json:"redis,optional"`

	ActiveTokenExpire       int `json:"activeTokenExpire,optional"`
	SessionTokenExpire      int `json:"sessionTokenExpire,optional"`
	DefaultRemainingTries   int `json:"defaultRemainingTries,optional"`
	DefaultRemainingAnswers int `json:"defaultRemainingAnswers,optional"`
}

type RestConf struct {
	rest.RestConf
}

type LogConf struct {
	logx.LogConf
}

type MySQLConf struct {
	DBName   string `json:"db_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

type BannerConf struct {
	Text     string `json:",default=JZERO"`
	Color    string `json:",default=green"`
	FontName string `json:",default=starwars,options=big|larry3d|starwars|standard"`
}
