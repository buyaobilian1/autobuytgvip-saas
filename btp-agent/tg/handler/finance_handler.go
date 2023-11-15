package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/pkg/id"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/service"
	tele "gopkg.in/telebot.v3"
)

func FinanceHandler(ctx tele.Context) error {
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}
	inlineBtn := &tele.ReplyMarkup{}

	inlineBtn.Inline(
		inlineBtn.Row(WithDrawBtn),
		inlineBtn.Row(CloseBtn, SupportBtn),
	)

	textFormat := `
尊敬的代理，您好：

钱包余额：%s USDT
钱包地址：%s

`
	yue := EscapeText(tele.ModeMarkdownV2, Float64Format(dbUser.Brokerage))
	addr := "未设置❌"
	if dbUser.TronAddr != nil && strings.TrimSpace(*dbUser.TronAddr) != "" {
		addr = *dbUser.TronAddr
	}
	text := fmt.Sprintf(textFormat, yue, addr)
	return ctx.Send(text, inlineBtn)
}

func SetTronAddrHandler(ctx tele.Context) error {
	//todo 波场地址检测
	var u = query.User
	addr := ctx.Data()
	if addr == "" {
		return ctx.Send(EscapeText(tele.ModeMarkdownV2, "请输入正确的USDT（TRC-20）提现地址"))
	}
	_, err := u.Where(u.TgID.Eq(ctx.Sender().ID)).Update(u.TronAddr, addr)
	if err != nil {
		log.Printf("update fail, %v\n", err)
		return err
	}
	text := "提现地址设置成功，当前地址 `" + addr + "`"
	return ctx.Send(text)
}

var NeedTronAddrText = `
请先设置提现钱包地址

相关指令
设置提现地址
/address T88888\*\*\*\*\*88888
`

func WithdrawHandler(ctx tele.Context) error {
	dbAgentUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}

	if dbAgentUser.TronAddr == nil || strings.TrimSpace(*dbAgentUser.TronAddr) == "" {
		return ctx.Send(NeedTronAddrText)
	} else if dbAgentUser.Brokerage <= 1 {
		return ctx.Send("代理账户余额不足")
	}

	return createWithdrawOrder(ctx)
}

func createWithdrawOrder(ctx tele.Context) error {
	var u = query.User
	dbAgentUser, err := u.Where(u.TgID.Eq(ctx.Sender().ID)).First()
	if err != nil {
		log.Printf("[db] query fail. %v \n", err)
		return err
	}

	withdrawModel := &model.Withdraw{
		UserID:    dbAgentUser.ID,
		TgID:      dbAgentUser.TgID,
		OrderNo:   id.GenerateId(1),
		Amount:    dbAgentUser.Brokerage,
		TronAddr:  *dbAgentUser.TronAddr,
		Status:    1,
		CreatedAt: time.Time{},
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		_, e := tx.User.Where(tx.User.ID.Eq(dbAgentUser.ID)).Update(tx.User.Brokerage, 0)
		if e != nil {
			return e
		}
		e = tx.Withdraw.Create(withdrawModel)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Printf("[db] udpate db fail. %v \n", err)
		return err
	}

	reply := "提现申请已提交，请等待系统处理"
	return ctx.Send(reply)
}
