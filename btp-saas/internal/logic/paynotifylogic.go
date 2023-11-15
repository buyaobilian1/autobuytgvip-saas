package logic

import (
	"context"
	"errors"
	"log"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/types"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq/handle"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/epusdt"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/proxy"
	"github.com/zeromicro/go-zero/core/logx"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type PayNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayNotifyLogic {
	return &PayNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayNotifyLogic) PayNotify(req *types.PayNotifyRequest) error {
	var u, r, o = query.User, query.Recharge, query.Order
	if req.Type == "order" { //订单支付成功回调
		log.Printf("[order] 收到回调消息：%+v\n", *req)
		sign, err := epusdt.Sign(req.PayNotifyBody, l.svcCtx.Config.PayConf.ApiToken)
		if err != nil {
			return err
		}
		if sign != req.Signature {
			return errors.New("sign fail")
		}
		_, err = o.Where(o.OrderNo.Eq(req.OrderId), o.Status.Eq(1)).Update(o.Status, req.Status)
		if err != nil {
			return err
		}

		task, _ := handle.NewGiftTelegramPremiumTask(req.OrderId)
		_, _ = mq.QueueClient.Enqueue(task)
		return nil

	} else if req.Type == "recharge" { //充值成功回调
		log.Printf("[recharge] 收到回调消息：%+v\n", *req)
		sign, err := epusdt.Sign(req.PayNotifyBody, l.svcCtx.Config.PayConf.ApiToken)
		if err != nil {
			return err
		}
		if sign != req.Signature {
			return errors.New("sign fail")
		}
		// 判断订单是不是存在
		dbRecharge, err := r.Where(r.OrderNo.Eq(req.OrderId), r.Status.Eq(1)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		err = query.Q.Transaction(func(tx *query.Query) error {
			if _, e := tx.Recharge.Where(tx.Recharge.OrderNo.Eq(req.OrderId)).Update(tx.Recharge.Status, 2); e != nil {
				return e
			}
			if recharge, e := tx.Recharge.Where(tx.Recharge.OrderNo.Eq(req.OrderId)).First(); e != nil {
				return e
			} else {
				if _, e = tx.User.Where(tx.User.ID.Eq(recharge.UserID)).Update(tx.User.Balance, tx.User.Balance.Add(recharge.ActualAmount)); e != nil {
					return e
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[recharge] update data fail. %v\n", err)
			return err
		}

		// 机器人通知充值成功
		res, err := u.Select(u.BotToken).Where(u.BotID.Eq(dbRecharge.BotID)).First()
		if err != nil {
			log.Printf("[db] query fail. %v \n", err)
			return nil
		}
		opt := tele.Settings{
			Token:   *res.BotToken,
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
		_, _ = tgBot.Send(user, "充值成功！！")
	}
	return nil
}
