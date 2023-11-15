package handler

import (
	tele "gopkg.in/telebot.v3"
)

func AgentKeyHandler(ctx tele.Context) error {
	return ctx.Send("代理未开启")
}
