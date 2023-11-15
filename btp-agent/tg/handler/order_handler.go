package handler

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/service"
	tele "gopkg.in/telebot.v3"
)

const PageSize = 5

var noOrderText = `
代理订单查询

暂无订单！
`

func OrderHandler(ctx tele.Context) error {
	_, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("find or create user fail. %v\n", err)
		return err
	}
	inlineBtn := &tele.ReplyMarkup{}
	threeMonthBtn := inlineBtn.Data("刷新订单", "refresh")
	inlineBtn.Inline(
		inlineBtn.Row(threeMonthBtn),
	)

	return agentOrderPageView(ctx, 1, false)
}

func AgentOrderPageHandler(ctx tele.Context, page int) error {
	return agentOrderPageView(ctx, page, true)
}

func agentOrderPageView(ctx tele.Context, page int, isEdit bool) error {
	rows, count, err := getAgentOrderPage(ctx.Sender().ID, page)
	if err != nil {
		log.Printf("[db] query page fail. %v\n", err)
		return err
	}

	inlineObj := &tele.ReplyMarkup{}

	if count > 0 {
		var prevBtn, nextBtn, counterBtn tele.Btn
		var inlineRows = make([]tele.Row, PageSize)
		for _, row := range rows {
			btnName := fmt.Sprintf("%s - (%.0f USDT)", row.OrderNo, row.UsdtAmount-row.BaseAmount)

			btn := inlineObj.Data(btnName, OrderDetailBtn, OrderDetailBtn, strconv.Itoa(int(row.ID)), strconv.Itoa(page))
			inlineRows = append(inlineRows, inlineObj.Row(btn))
		}

		if page <= 1 {
			prevBtn = inlineObj.Data("✖️", "none", "1")
		} else {
			prevBtn = inlineObj.Data("⬅️", OrderPagePrevBtnId, OrderPagePrevBtnId, strconv.Itoa(page-1))
		}

		totalCountFloat64 := float64(count) / float64(PageSize)
		totalCount := int(math.Ceil(totalCountFloat64))
		counterBtn = inlineObj.Data(fmt.Sprintf("%d/%d", page, totalCount), "counter", "1")

		if page >= totalCount {
			nextBtn = inlineObj.Data("✖️", "none", "1")
		} else {
			nextBtn = inlineObj.Data("➡️", OrderPagePrevBtnId, OrderPagePrevBtnId, strconv.Itoa(page+1))
		}

		inlineRows = append(inlineRows, inlineObj.Row(prevBtn, counterBtn, nextBtn))

		inlineObj.Inline(
			inlineRows...,
		)
	}

	if count <= 0 {
		inlineObj.Inline(
			inlineObj.Row(CloseBtn, SupportBtn),
		)
		return ctx.Send(noOrderText, inlineObj)
	}

	text := "代理订单查询！"
	if isEdit {
		return ctx.Edit(text, inlineObj)
	}
	return ctx.Send(text, inlineObj)
}

func ShowAgentOrderDetail(ctx tele.Context) error {
	var o = query.Order
	split := strings.Split(ctx.Data(), "|")
	orderId, _ := strconv.Atoi(split[1])
	page, _ := strconv.Atoi(split[2])
	dbOrder, err := o.Where(o.ID.Eq(int32(orderId))).First()
	if err != nil {
		log.Printf("[db] query data fail. %v\n", err)
		return err
	}

	replyFormat := `
代理订单详情

订单：%s
商品：%s
价格：%s
佣金：%s
时间：%s

`
	orderNo := fmt.Sprintf("`%s`", dbOrder.OrderNo)
	goods := fmt.Sprintf("%d个月TG会员", dbOrder.VipMonth)
	price := fmt.Sprintf("%.0f USDT", dbOrder.UsdtAmount)
	brokerage := fmt.Sprintf("%.0f USDT", dbOrder.UsdtAmount-dbOrder.BaseAmount)
	ctime := EscapeText(tele.ModeMarkdownV2, dbOrder.CreatedAt.Format(time.DateTime))
	reply := fmt.Sprintf(replyFormat, orderNo, goods, price, brokerage, ctime)

	inlineObj := &tele.ReplyMarkup{}
	backBtn := inlineObj.Data("⬅️返回", OrderDetailBackBtn, OrderDetailBackBtn, strconv.Itoa(page))
	inlineObj.Inline(
		inlineObj.Row(backBtn),
	)

	return ctx.Edit(reply, inlineObj)
}

// 分页查询代理的订单
func getAgentOrderPage(agentTgId int64, page int) (result []*model.Order, count int64, err error) {
	var o, u = query.Order, query.User
	offset := (page - 1) * PageSize
	result, count, err = o.Select(o.ALL).LeftJoin(u, u.ID.EqCol(o.AgentUserID)).
		Where(u.TgID.Eq(agentTgId), o.Status.Eq(3), o.AgentStatus.Eq(2)).FindByPage(offset, PageSize)
	return
}
