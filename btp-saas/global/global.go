package global

import (
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/config"
	tele "gopkg.in/telebot.v3"
)

var (
	Conf      config.Config
	BotMapper = make(map[int64]*tele.Bot)
)
