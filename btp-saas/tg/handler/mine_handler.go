package handler

import (
	"fmt"
	"log"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/service"
	tele "gopkg.in/telebot.v3"
)

func MineKeyHandler(ctx tele.Context) error {
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("[db] query fila. : %v, dbuser: %v", err, dbUser)
		return err
	}

	replyFormat := `
个人信息

用户昵称：@%s
钱包余额：%s USDT
`
	balanceStr := EscapeText(tele.ModeMarkdownV2, fmt.Sprintf("%.2f", dbUser.Balance))
	text := fmt.Sprintf(replyFormat, ctx.Sender().Username, balanceStr)

	replyMarkup := &tele.ReplyMarkup{}
	replyMarkup.Inline(
		replyMarkup.Row(RechargeBtn),
		replyMarkup.Row(CloseBtn, SupportBtn),
	)

	return ctx.Send(text, replyMarkup)
}
