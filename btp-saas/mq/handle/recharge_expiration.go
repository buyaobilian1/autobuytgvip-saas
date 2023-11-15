package handle

import (
	"context"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/proxy"
	"github.com/hibiken/asynq"
	tele "gopkg.in/telebot.v3"
)

const RechargeExpirationPattern = "recharge:expiration"

func NewRechargeExpirationTask(orderNo string) (*asynq.Task, error) {
	return asynq.NewTask(RechargeExpirationPattern, []byte(orderNo)), nil
}

func RechargeExpirationHandler(ctx context.Context, t *asynq.Task) error {
	var r, u = query.Recharge, query.User
	orderNo := string(t.Payload())
	dbRecharge, err := r.Where(r.OrderNo.Eq(orderNo)).First()
	if err != nil {
		return err
	}
	if dbRecharge.Status == 1 {
		_, err = r.Where(r.OrderNo.Eq(orderNo)).Update(r.Status, 3)
	} else {
		return nil
	}

	// 通知机器人
	type Result struct {
		BotToken string
	}
	res := Result{}
	_ = u.Select(u.BotToken).LeftJoin(r, u.BotID.EqCol(r.BotID)).Where(r.OrderNo.Eq(orderNo)).Scan(&res)
	opt := tele.Settings{
		Token:   res.BotToken,
		Offline: true,
	}
	if global.Conf.AppConf.ProxyUrl != "" {
		opt.Client = proxy.NewProxyHttpClient(global.Conf.AppConf.ProxyUrl)
	}
	tgBot, _ := tele.NewBot(opt)
	_ = tgBot.Delete(&tele.Message{
		ID: int(dbRecharge.TgMsgID),
		Chat: &tele.Chat{
			ID: dbRecharge.TgChatID,
		},
	})
	user := &tele.User{
		ID: dbRecharge.TgChatID,
	}
	_, _ = tgBot.Send(user, "🚫支付超时，订单已取消")
	return nil
}
