package handler

import (
	tele "gopkg.in/telebot.v3"
)

func CooperationKeyHandler(ctx tele.Context) error {
	return SupportHandler(ctx)
}
