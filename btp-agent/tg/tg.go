package tg

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/config"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/tg/handler"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func Start(c config.AppConf) {
	bot, err := createBot(c)
	if err != nil {
		log.Fatalf("agent bot startup fail. %v\n", err)
	}

	initSimple(bot)
	go bot.Start()
}

func createBot(c config.AppConf) (*tele.Bot, error) {
	settings := tele.Settings{
		Token:     c.BotToken,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeMarkdownV2,
	}
	if c.ProxyUrl != "" {
		settings.Client = global.NewProxyHttpClient(c.ProxyUrl)
	}
	return tele.NewBot(settings)
}

func initSimple(b *tele.Bot) {
	startCmd := []tele.Command{
		{Text: "/start", Description: "开始"},
		{Text: "/help", Description: "帮助"},
		{Text: "/token", Description: "设置机器人Token"},
		{Text: "/price3", Description: "设置3个月会员销售价"},
		{Text: "/price6", Description: "设置6个月会员销售价"},
		{Text: "/price12", Description: "设置12个月会员销售价"},
		{Text: "/orders", Description: "代理订单查询"},
		{Text: "/address", Description: "设置提现钱包地址"},
		{Text: "/withdraw", Description: "发起提现"},
	}
	_ = b.SetCommands(startCmd)
	b.Use(middleware.AutoRespond())
	b.Handle("/start", handler.StartHandler)
	b.Handle("/help", handler.HelpHandler)
	b.Handle("/token", handler.BotTokenHandler)
	b.Handle("/price3", handler.Price3Handler)
	b.Handle("/price6", handler.Price6Handler)
	b.Handle("/price12", handler.Price12Handler)
	b.Handle("/orders", handler.OrderHandler)
	b.Handle("/address", handler.SetTronAddrHandler)
	b.Handle("/withdraw", handler.WithdrawHandler)
	// inline button
	b.Handle(&handler.CloseBtn, handler.CloseHandler)
	b.Handle(&handler.SupportBtn, handler.SupportHandler)
	b.Handle(&handler.BotTokenSettingBtn, handler.BotTokenSettingHandler)
	b.Handle(&handler.AgentPriceSettingBtn, handler.PriceHandler)
	b.Handle(&handler.WithDrawBtn, handler.WithdrawHandler)
	// keyboard
	b.Handle(&handler.FinanceKeyboard, handler.FinanceHandler)
	b.Handle(&handler.OrderKeyboard, handler.OrderHandler)
	b.Handle(&handler.BotKeyboard, handler.BotHandler)

	b.Handle(tele.OnCallback, func(ctx tele.Context) error {
		dataArr := strings.Split(ctx.Data(), "|")

		if dataArr[0] == handler.OrderPagePrevBtnId || dataArr[0] == handler.OrderPageNextBtnId {
			page, _ := strconv.Atoi(dataArr[1])
			return handler.AgentOrderPageHandler(ctx, page)
		} else if dataArr[0] == handler.OrderDetailBtn {
			return handler.ShowAgentOrderDetail(ctx)
		} else if dataArr[0] == "\f"+handler.OrderDetailBackBtn {
			page, _ := strconv.Atoi(dataArr[2])
			return handler.AgentOrderPageHandler(ctx, page)
		}

		return nil
	})

}
