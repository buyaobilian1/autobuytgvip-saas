package handler

import (
	"fmt"
	"log"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/service"
	tele "gopkg.in/telebot.v3"
)

var unConfigFormat = `
尊敬的用户，你好：

当前代理尚未开通成功！！
👇请按照以下顺序进行设置👇

%s

👇点击下方【⚙️设置】👇
`

var configOkFormat = `
尊敬的用户 @%s，欢迎回来！
`

func StartHandler(ctx tele.Context) error {
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}

	keyboards := &tele.ReplyMarkup{ResizeKeyboard: true}
	keyboards.Reply(
		keyboards.Row(BotKeyboard, OrderKeyboard, FinanceKeyboard),
	)

	// 代理未配置完成的情况
	if dbUser.BotToken == nil || dbUser.BotID == nil || dbUser.ThreeMonthPrice == nil || dbUser.SixMonthPrice == nil || dbUser.TwelveMonthPrice == nil {
		checkText := getAgentConfigCheckInfo(dbUser)
		reply := fmt.Sprintf(unConfigFormat, checkText)
		return ctx.Send(reply, keyboards)
	}
	// 已正常配置的代理用户
	replay := fmt.Sprintf(configOkFormat, ctx.Sender().Username)
	return ctx.Send(replay, keyboards)
}

var checkFomater = `1️⃣ 机器人Token %s
2️⃣  3个月会员价 %s
3️⃣  6个月会员价 %s
4️⃣12个月会员价 %s`

func getAgentConfigCheckInfo(dbUser *model.User) string {
	var tokenStr, price3, price6, price12 = "未设置❌", "未设置❌", "未设置❌", "未设置❌"
	if dbUser.BotToken != nil {
		tokenStr = "已设置✅"
	}
	if dbUser.ThreeMonthPrice != nil {
		price3 = fmt.Sprintf("%s USDT ✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.ThreeMonthPrice)))
	}
	if dbUser.SixMonthPrice != nil {
		price6 = fmt.Sprintf("%s USDT ✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.SixMonthPrice)))
	}
	if dbUser.TwelveMonthPrice != nil {
		price12 = fmt.Sprintf("%s USDT ✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.TwelveMonthPrice)))
	}

	return fmt.Sprintf(checkFomater, tokenStr, price3, price6, price12)
}
