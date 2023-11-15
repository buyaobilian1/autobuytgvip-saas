package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConf MysqlConf
	RedisConf RedisConf
	PayConf   PayConf
	AppConf   AppConf
}

type MysqlConf struct {
	Host string
	User string
	Pass string
	Db   string
}

type RedisConf struct {
	Host string
	Type string
	Pass string
	Tls  bool
}

type PayConf struct {
	BaseApi   string
	ApiToken  string
	NotifyUrl string
}

type AppConf struct {
	OrderExpireMinute    int
	RechargeExpireMinute int
	ProxyUrl             string
	TonSeed              string
	Hash                 string
	Cookie               string
}
