package handler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/service"
	tele "gopkg.in/telebot.v3"
)

func BotHandler(ctx tele.Context) error {
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}
	inlineBtn := &tele.ReplyMarkup{}

	inlineBtn.Inline(
		inlineBtn.Row(BotTokenSettingBtn),
		inlineBtn.Row(AgentPriceSettingBtn),
		inlineBtn.Row(CloseBtn, SupportBtn),
	)

	textFormat := `
⭐️⭐️⭐️注意注意注意⭐️⭐️⭐️

Token设置后，请不要随意更换，否则机器人功能将失效

状态信息
%s

`
	checkText := getAgentConfigCheckInfo(dbUser)
	reply := fmt.Sprintf(textFormat, checkText)
	return ctx.Send(reply, inlineBtn)
}

func BotTokenSettingHandler(ctx tele.Context) error {
	text := `
机器人设置教程：
输入指令 /token BotFather中复制过来的token
例如 我的token是 123456:abcdef
则输入 /token 123456:abcdef
`
	return ctx.Send(text)
}

func BotTokenHandler(ctx tele.Context) error {
	var u = query.User
	token := ctx.Data()
	if ok, _ := regexp.MatchString("^\\d{10}:[a-zA-Z0-9_-]{35}", token); !ok {
		return ctx.Send("token无效")
	}
	split := strings.Split(token, ":")
	botId, err := strconv.Atoi(split[0])
	if err != nil {
		return ctx.Send("token无效")
	}
	bot, err := checkBotToken(token)
	if err != nil {
		return ctx.Send("token无效")
	}
	_ = ctx.Send("系统将为你自动检测和配置机器人，请耐心等待几秒钟!!")

	// 入库
	botIdInt64 := int64(botId)
	now := time.Now()
	updateUser := model.User{
		BotID:        &botIdInt64,
		BotToken:     &token,
		BotUsername:  &bot.Me.Username,
		BotStatus:    1,
		BotCreatedAt: &now,
	}
	_, _ = u.Where(u.TgID.Eq(ctx.Sender().ID)).Updates(updateUser)

	webhook := &tele.Webhook{
		SecretToken: global.Conf.AppConf.WebhookSecret,
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: fmt.Sprintf(global.Conf.AppConf.WebhookUrl, botIdInt64),
		},
	}
	err = bot.SetWebhook(webhook)
	if err != nil {
		return ctx.Send("机器人配置失败")
	}

	return ctx.Send("机器人Token设置成功")
}

func checkBotToken(token string) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:   token,
		Offline: false,
	}
	if global.Conf.AppConf.ProxyUrl != "" {
		pref.Client = global.NewProxyHttpClient(global.Conf.AppConf.ProxyUrl)
	}
	return tele.NewBot(pref)
}
