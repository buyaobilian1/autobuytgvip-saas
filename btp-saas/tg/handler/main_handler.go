package handler

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/fragment"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/id"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/image"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/service"
	tele "gopkg.in/telebot.v3"
)

func StartHandler(ctx tele.Context) error {
	_, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		return err
	}
	dbUser, err := service.FindUserByBotId(ctx.Bot().Me.ID)
	if err != nil {
		return err
	}
	var startFormatText = `
	❤️本机器人向您提供Telegram Premium会员自动开通服务！
	
	当前价格：
	*  3个月 / %s U*
	*  6个月 / %s U*
	*12个月 / %s U*🔥
	
	请选择下方菜单：`
	price3 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.ThreeMonthPrice))
	price6 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.SixMonthPrice))
	price12 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.TwelveMonthPrice))
	startText := fmt.Sprintf(startFormatText, price3, price6, price12)
	keyboards := &tele.ReplyMarkup{ResizeKeyboard: true}
	keyboards.Reply(
		keyboards.Row(RechargeKeyboard, MineKeyboard),
		keyboards.Row(AgentKeyboard, CooperationKeyboard),
	)
	_ = ctx.Send(startText, keyboards)

	startText2 := "尊贵的Telegram用户您好！\n\n请选择为谁开通/续费TG会员："
	replyMarkup := &tele.ReplyMarkup{}
	btnBuyMyself := replyMarkup.Data("✈️此账号开通", ByMyselfBtnId, "@"+ctx.Sender().Username)
	btnGift := replyMarkup.Data("🎁赠送给他人", GiftOtherBtnId)
	replyMarkup.Inline(
		replyMarkup.Row(btnBuyMyself, btnGift),
		replyMarkup.Row(CloseBtn, SupportBtn),
	)
	return ctx.Send(EscapeText(tele.ModeMarkdownV2, startText2), replyMarkup)
}

func BuyMyselfHandler(ctx tele.Context) error {
	arr := strings.Split(ctx.Data(), "|")
	return ShowTgUserInfo(ctx, arr[1])
}

func GiftOtherHandler(ctx tele.Context) error {
	var giftFormatText = `
请直接发送你需要开通会员的Telegram用户名：

*提示：用户名以@开头，如 %s*
`
	giftText := fmt.Sprintf(giftFormatText, "@"+ctx.Sender().Username)
	replyMarkup := &tele.ReplyMarkup{
		ForceReply:     true,
		Placeholder:    "请输入Tg用户名",
		ResizeKeyboard: true,
	}

	return ctx.Send(giftText, replyMarkup)
}

func BuyThreeMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func BuySixMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func BuyTwelveMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func CreateTelegramPremiumOrder(ctx tele.Context) error {
	var u, o, p = query.User, query.Order, query.Param
	dbAgentUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("[db] query data fail, %v\n", err)
		return err
	}
	dbUser, err := u.Where(u.TgID.Eq(ctx.Sender().ID)).First()
	if err != nil {
		log.Printf("[db] query data fail, %v\n", err)
		return err
	}
	basePriceObj, err := p.Where(p.K.Eq("base_price")).First()
	if err != nil {
		log.Printf("[db] query data fail. %v\n", err)
		return err
	}
	orderFormatText := `❗️❗️❗️请注意：网络必须是TRC\-20，否则无法到账
❗️❗️❗️请注意，金额必须与下面的一致（一位都不能少）
👇*请向以下地址转账 %s USDT*

%s

👆点击复制上面地址进行支付，或者扫描上面二维码支付。
`

	params := strings.Split(ctx.Data(), "|")
	distUsername := params[1]
	vipMonth, _ := strconv.Atoi(params[3])
	var usdtAmount, baseAmount float64
	if vipMonth == 3 {
		usdtAmount = *dbAgentUser.ThreeMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V1, 64)
		if err != nil {
			log.Printf("base price3 set fail. %v\n", err)
			return err
		}
	} else if vipMonth == 6 {
		usdtAmount = *dbAgentUser.SixMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V2, 64)
		if err != nil {
			log.Printf("base price3 set fail. %v\n", err)
			return err
		}
	} else if vipMonth == 12 {
		usdtAmount = *dbAgentUser.TwelveMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V3, 64)
		if err != nil {
			log.Printf("base price3 set fail. %v\n", err)
			return err
		}
	} else {
		return errors.New("套餐错误")
	}

	order := &model.Order{
		OrderNo:           id.GenerateId(1),
		UserID:            dbUser.ID,
		AgentUserID:       dbAgentUser.ID,
		BotID:             ctx.Bot().Me.ID,
		ReceiveTgUsername: distUsername,
		VipMonth:          int32(vipMonth),
		UsdtAmount:        usdtAmount,
		BaseAmount:        baseAmount,
		Status:            1,
		CreatedAt:         time.Now(),
		ExpiredAt:         time.Now().Add(10 * time.Minute),
		TgChatID:          ctx.Chat().ID,
		TgMsgID:           0, //先置零，登消息发出去后得到消息id后再更新
	}
	res, err := service.CreateOrder(order)
	if err != nil {
		log.Printf("fail to create order: %v\n", err)
		return ctx.Respond(&tele.CallbackResponse{
			Text:      "系统繁忙，订单创建失败，请重试",
			ShowAlert: true,
		})
	}

	replyMarkup := &tele.ReplyMarkup{}
	replyMarkup.Inline(
		replyMarkup.Row(SupportBtn),
	)

	amountStr := Float64Format(res.ActualAmount)
	replyText := fmt.Sprintf(orderFormatText, EscapeText(tele.ModeMarkdownV2, amountStr), res.Token)
	context := &tele.Photo{
		File:    tele.FromReader(image.GenQrcode(res.Token)),
		Caption: replyText,
	}
	msg, err := ctx.Bot().Send(ctx.Recipient(), context, replyMarkup)

	_, err = o.Where(o.OrderNo.Eq(order.OrderNo)).Update(o.TgMsgID, msg.ID)
	return err
}

func ShowTgUserInfo(ctx tele.Context, username string) error {
	var u = query.User
	dbUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("query db fail: %v\n", err)
		return err
	}
	var replyText string
	var replyFormatText = `
开通用户：%s
用户昵称：%s

确定为此用户 开通/续费 Telegram Premium会员吗？
`
	userInfo, err := fragment.SearchPremiumGiftRecipient(username, 3)
	if err != nil {
		log.Printf("fail to get premium gift recipient, %v", err)
		return nil
	}
	if userInfo.Error == "No Telegram users found." {
		return ctx.Send(EscapeText(tele.ModeMarkdownV2, "用户名不存在."))
	}
	if userInfo.Error == "This account is already subscribed to Telegram Premium." {
		return ctx.Send(EscapeText(tele.ModeMarkdownV2, "此账号已经订阅会员."))
	}

	replyText = EscapeText(tele.ModeMarkdownV2, fmt.Sprintf(replyFormatText, username, userInfo.Found.Name))

	replyMarkup := &tele.ReplyMarkup{}
	btnBuy3Month := replyMarkup.Data(fmt.Sprintf("3个月 / %s U", Float64Format(*dbUser.ThreeMonthPrice)), BuyThreeMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "3")
	btnBuy6Month := replyMarkup.Data(fmt.Sprintf("6个月 / %s U", Float64Format(*dbUser.SixMonthPrice)), BuySixMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "6")
	btnBuy12Month := replyMarkup.Data(fmt.Sprintf("12个月 / %s U🔥", Float64Format(*dbUser.TwelveMonthPrice)), BuyTwelveMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "12")

	replyMarkup.Inline(
		replyMarkup.Row(btnBuy3Month, btnBuy6Month),
		replyMarkup.Row(btnBuy12Month),
		replyMarkup.Row(CloseBtn, SupportBtn),
	)

	return ctx.Send(replyText, replyMarkup)
}
