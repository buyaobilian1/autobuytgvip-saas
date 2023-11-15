package tg

import (
	"errors"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/proxy"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/tg/handler"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"regexp"
	"strings"
)

func GetBotInstance(botId int64) (*tele.Bot, error) {
	if v, ok := global.BotMapper[botId]; ok {
		return v, nil
	}
	var u = query.User
	agentUser, err := u.Where(u.BotStatus.Eq(1), u.BotID.Eq(botId), u.BotToken.IsNotNull(),
		u.ThreeMonthPrice.IsNotNull(), u.SixMonthPrice.IsNotNull(), u.TwelveMonthPrice.IsNotNull()).First()
	if err != nil {
		return nil, err
	}
	bot, err := NewSingleBot(agentUser)
	if err != nil {
		return nil, err
	}
	global.BotMapper[botId] = bot

	return bot, nil
}

func StartAll() {
	var u = query.User
	userList, err := u.Where(u.BotStatus.Eq(1), u.BotToken.IsNotNull()).Find()
	if err != nil {
		panic(err)
	}
	if len(userList) <= 0 {
		log.Printf("暂无可运行的机器人.\n")
	}
	for _, user := range userList {
		bot, err := NewSingleBot(user)
		if err != nil {
			continue
		}
		global.BotMapper[bot.Me.ID] = bot
	}

}

func NewSingleBot(user *model.User) (*tele.Bot, error) {
	if user.BotID == nil || user.BotToken == nil || user.ThreeMonthPrice == nil || user.SixMonthPrice == nil || user.TwelveMonthPrice == nil {
		err := errors.New("bot config missing")
		log.Printf("[bot %d] startup fail.%v\n", user.BotID, err)
		return nil, err
	}
	pref := tele.Settings{
		Token:     *user.BotToken,
		ParseMode: tele.ModeMarkdownV2,
		OnError: func(err error, ctx tele.Context) {
			log.Printf("[bot] message err %v\n", err)
		},
	}
	if global.Conf.AppConf.ProxyUrl != "" {
		pref.Client = proxy.NewProxyHttpClient(global.Conf.AppConf.ProxyUrl)
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Printf("[bot %d] start fail.%v\n", user.BotID, err)
		return nil, err
	}

	setCommands(bot)
	setBiz(bot)

	return bot, nil

}

func setCommands(bot *tele.Bot) {
	startCmd := tele.Command{
		Text:        "/start",
		Description: "开始下单",
	}
	err := bot.SetCommands(startCmd)
	if err != nil {
		log.Printf("[bot %d] set command fail. %v\n", bot.Me.ID, err)
	}
}

func setBiz(bot *tele.Bot) {
	bot.Use(middleware.AutoRespond())
	bot.Handle("/start", handler.StartHandler)
	bot.Handle(tele.OnCallback, func(ctx tele.Context) error {
		dataArr := strings.Split(ctx.Data(), "|")
		matchStr := strings.ReplaceAll(dataArr[0], "\f", "")
		if matchStr == handler.ByMyselfBtnId {
			return handler.BuyMyselfHandler(ctx)
		} else if matchStr == handler.GiftOtherBtnId {
			return handler.GiftOtherHandler(ctx)
		} else if matchStr == handler.CloseBtnId {
			return handler.CloseHandler(ctx)
		} else if matchStr == handler.SupportBtnId {
			return handler.SupportHandler(ctx)
		} else if matchStr == handler.RechargeBtnId {
			return handler.RechargeKeyHandler(ctx)
		} else if matchStr == handler.BuyThreeMonthBtnId {
			return handler.BuyThreeMonthHandler(ctx)
		} else if matchStr == handler.BuySixMonthBtnId {
			return handler.BuySixMonthHandler(ctx)
		} else if matchStr == handler.BuyTwelveMonthBtnId {
			return handler.BuyTwelveMonthHandler(ctx)
		} else if matchStr == handler.RechargeConfirmBtnId {
			return handler.RechargeDoHandler(ctx)
		}

		log.Printf("callback data => %v\n", dataArr)

		return ctx.Respond()
	})
	bot.Handle(tele.OnText, func(ctx tele.Context) error {
		if ok, _ := regexp.MatchString("^@.*", ctx.Text()); ok {
			return handler.ShowTgUserInfo(ctx, ctx.Text())
		} else if ok, _ = regexp.MatchString("^[123456789]\\d*$", ctx.Text()); ok {
			return handler.RechargeConfirm(ctx)
		}

		return nil
	})

	// 键盘处理
	bot.Handle(&handler.RechargeKeyboard, handler.RechargeKeyHandler)
	bot.Handle(&handler.MineKeyboard, handler.MineKeyHandler)
	bot.Handle(&handler.AgentKeyboard, handler.AgentKeyHandler)
	bot.Handle(&handler.CooperationKeyboard, handler.CooperationKeyHandler)
}
