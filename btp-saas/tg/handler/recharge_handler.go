package handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/id"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/image"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/service"
	tele "gopkg.in/telebot.v3"
)

func RechargeKeyHandler(ctx tele.Context) error {
	text := `
您好，欢迎使用充值服务：

*充值金额为整数*
*最小充值 1 U*

请直接输入你需要充值的金额：
`
	return ctx.Send(text)
}

func RechargeConfirm(ctx tele.Context) error {
	replyFormat := "您确定需要充值 %s USDT吗？"

	inlineObj := &tele.ReplyMarkup{}
	okBtn := inlineObj.Data("十分确定", RechargeConfirmBtnId, ctx.Text())
	inlineObj.Inline(
		inlineObj.Row(okBtn, CloseBtn),
	)
	reply := fmt.Sprintf(replyFormat, ctx.Text())
	return ctx.Send(reply, inlineObj)
}

func RechargeDoHandler(ctx tele.Context) error {
	arrs := strings.Split(ctx.Data(), "|")
	var u = query.User
	dbUser, err := u.Where(u.TgID.Eq(ctx.Sender().ID)).First()
	if err != nil {
		log.Printf("db fail.%v\n", err)
		return err
	}
	format := `❗️❗️❗️请注意：网络必须是TRC\-20，否则无法到账
❗️❗️❗️请注意，金额必须与下面的一致（一位都不能少）
👇*请向以下地址转账 %s USDT*

%s

👆点击复制上面地址进行支付，或者扫描上面二维码支付。
`
	amount, _ := strconv.Atoi(arrs[1])

	v := &model.Recharge{
		UserID:       dbUser.ID,
		BotID:        ctx.Bot().Me.ID,
		OrderNo:      id.GenerateId(1),
		Amount:       float64(amount),
		Status:       1,
		ActualAmount: 0,
		CreatedAt:    time.Time{},
		TgChatID:     ctx.Chat().ID,
		TgMsgID:      0,
	}
	result, err := service.CreateRechargeOrder(v)
	if err != nil {
		log.Printf("fail to create order: %v\n", err)
		return ctx.Respond(&tele.CallbackResponse{
			Text:      "系统繁忙，订单创建失败，请重试",
			ShowAlert: true,
		})
	}
	amountStr := EscapeText(tele.ModeMarkdownV2, Float64Format(result.ActualAmount))
	replyText := fmt.Sprintf(format, amountStr, result.Token)
	photo := &tele.Photo{
		File:    tele.FromReader(image.GenQrcode(result.Token)),
		Caption: replyText,
	}

	msg, err := ctx.Bot().Send(ctx.Recipient(), photo)
	var rechargeDao = query.Recharge
	_, err = rechargeDao.Where(rechargeDao.OrderNo.Eq(v.OrderNo)).Update(rechargeDao.TgMsgID, msg.ID)

	return err
}
