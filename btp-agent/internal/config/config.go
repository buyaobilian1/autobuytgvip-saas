package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MysqlConf MysqlConf
	AppConf   AppConf
}

type MysqlConf struct {
	Host string
	User string
	Pass string
	Db   string
}

type AppConf struct {
	BotToken      string
	ProxyUrl      string
	WebhookUrl    string
	WebhookSecret string
}
