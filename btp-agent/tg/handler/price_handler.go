package handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/service"
	tele "gopkg.in/telebot.v3"
)

func PriceHandler(ctx tele.Context) error {
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}
	var p = query.Param
	basePriceParam, err := p.Where(p.K.Eq("base_price")).First()
	if err != nil {
		return err
	}

	textFormat := `
价格设置必须高于代理价，不然没有利润哦。

当前代理价
  3个月会员 %s USDT
  6个月会员 %s USDT
12个月会员 %s USDT

当前销售价
%s
%s
%s

可使用如下指令设置销售价：
/price3 3个月的价格 
/price6 6个月的价格 
/price12 12个月的价格
⭐️⭐️注意⭐️⭐️
⭐️⭐️价格只能设置整数⭐️⭐️
如：你想设置3个月会员价为15U
则：/price3 15

`
	basePrice3Str := EscapeText(tele.ModeMarkdownV2, *basePriceParam.V1)
	basePrice6Str := EscapeText(tele.ModeMarkdownV2, *basePriceParam.V2)
	basePrice12Str := EscapeText(tele.ModeMarkdownV2, *basePriceParam.V3)

	var price3Str, price6Str, price12Str = "  3个月会员 未设置❌", "  6个月会员 未设置❌", "12个月会员 未设置❌"
	if dbUser.ThreeMonthPrice != nil {
		price3Str = fmt.Sprintf("  3个月会员 %s USDT✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.ThreeMonthPrice)))
	}
	if dbUser.SixMonthPrice != nil {
		price6Str = fmt.Sprintf("  6个月会员 %s USDT✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.SixMonthPrice)))
	}
	if dbUser.TwelveMonthPrice != nil {
		price12Str = fmt.Sprintf("12个月会员 %s USDT✅", EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.TwelveMonthPrice)))
	}

	text := fmt.Sprintf(textFormat, basePrice3Str, basePrice6Str, basePrice12Str, price3Str, price6Str, price12Str)

	return ctx.Send(text)
}

func Price3Handler(ctx tele.Context) error {
	return setPrice(ctx, 3)
}

func Price6Handler(ctx tele.Context) error {
	return setPrice(ctx, 6)
}

func Price12Handler(ctx tele.Context) error {
	return setPrice(ctx, 12)
}

func setPrice(ctx tele.Context, month int) error {
	var p, u = query.Param, query.User

	priceIn, err := strconv.Atoi(strings.TrimSpace(ctx.Data()))
	if err != nil {
		return ctx.Send("价格必须输入整数！")
	}
	fmt.Println(priceIn)

	dbUser, _ := p.Where(p.K.Eq("base_price")).First()

	var dbPrice float64
	if month == 3 {
		dbPrice, _ = strconv.ParseFloat(*dbUser.V1, 64)
	} else if month == 6 {
		dbPrice, _ = strconv.ParseFloat(*dbUser.V2, 64)
	} else if month == 12 {
		dbPrice, _ = strconv.ParseFloat(*dbUser.V3, 64)
	}

	if float64(priceIn) < dbPrice {
		return ctx.Send("销售价不能低于代理价哦。不然裤衩子要赔进去哦！")
	}

	if month == 3 {
		_, err = u.Where(u.TgID.Eq(ctx.Sender().ID)).Update(u.ThreeMonthPrice, float64(priceIn))
	} else if month == 6 {
		_, err = u.Where(u.TgID.Eq(ctx.Sender().ID)).Update(u.SixMonthPrice, float64(priceIn))
	} else if month == 12 {
		_, err = u.Where(u.TgID.Eq(ctx.Sender().ID)).Update(u.TwelveMonthPrice, float64(priceIn))
	}

	return ctx.Send(fmt.Sprintf("%d个月会员的销售价格设置成功", month))
}
