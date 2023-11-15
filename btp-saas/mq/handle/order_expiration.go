package handle

import (
	"context"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/proxy"
	"github.com/hibiken/asynq"
	tele "gopkg.in/telebot.v3"
)

const OrderExpirationPattern = "order:expiration"

func NewOrderExpirationTask(orderNo string) (*asynq.Task, error) {
	return asynq.NewTask(OrderExpirationPattern, []byte(orderNo)), nil
}

func OrderExpirationHandler(ctx context.Context, t *asynq.Task) error {
	var o, u = query.Order, query.User
	orderNo := string(t.Payload())
	dbOrder, err := o.Where(o.OrderNo.Eq(orderNo)).First()
	if err != nil {
		return err
	}
	if dbOrder.Status == 1 {
		if _, err = o.Where(o.OrderNo.Eq(orderNo)).Update(o.Status, 4); err != nil {
			return err
		}
	} else {
		return nil
	}
	// 通知机器人
	type Result struct {
		BotToken string
	}
	res := Result{}
	_ = u.Select(u.BotToken).LeftJoin(o, u.BotID.EqCol(o.BotID)).Where(o.OrderNo.Eq(orderNo)).Scan(&res)
	opt := tele.Settings{
		Token:   res.BotToken,
		Offline: true,
	}
	if global.Conf.AppConf.ProxyUrl != "" {
		opt.Client = proxy.NewProxyHttpClient(global.Conf.AppConf.ProxyUrl)
	}
	tgBot, _ := tele.NewBot(opt)
	_ = tgBot.Delete(&tele.Message{
		ID: int(dbOrder.TgMsgID),
		Chat: &tele.Chat{
			ID: dbOrder.TgChatID,
		},
	})
	user := &tele.User{
		ID: dbOrder.TgChatID,
	}
	_, _ = tgBot.Send(user, "🚫支付超时，订单已取消")

	return nil
}
