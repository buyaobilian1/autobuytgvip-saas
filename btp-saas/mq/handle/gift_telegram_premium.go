package handle

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/blockchain"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/fragment"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/proxy"
	"github.com/hibiken/asynq"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

const GiftTelegramPremiumPattern = "premium:gift"

func NewGiftTelegramPremiumTask(orderNo string) (*asynq.Task, error) {
	return asynq.NewTask(GiftTelegramPremiumPattern, []byte(orderNo)), nil
}

func GiftTelegramPremiumHandler(ctx context.Context, t *asynq.Task) error {
	var o, u = query.Order, query.User
	orderNo := string(t.Payload())
	// 获取并判断订单信息
	dbOrder, err := o.Where(o.OrderNo.Eq(orderNo), o.Status.Eq(2)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("order is not found or status is not eq 2. orderNo = %s", orderNo)
			return nil
		}
		return err
	}

	// 自动开通爬虫
	fragmentRefId, err := buyTelegramPremium(dbOrder.ReceiveTgUsername, int(dbOrder.VipMonth))
	if err != nil {
		log.Printf("[crawler] buy telegram premium fail. %v\n", err)
		return nil
	}
	err = updateSuccess(orderNo, fragmentRefId)
	if err != nil {
		log.Printf("[db] 会员已开通，但是更新数据库失败. %v\n", err)
		return nil
	}

	dbUser, err := u.Where(u.ID.Eq(dbOrder.UserID)).First()
	if err != nil {
		return nil
	}

	opt := tele.Settings{
		Token:   *dbUser.BotToken,
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
	format := "恭喜您，成功开通 %d 个月Tg会员。"
	_, _ = tgBot.Send(user, fmt.Sprintf(format, dbOrder.VipMonth))

	return nil
}

var tonCommentFormats = map[int]string{
	3:  "Telegram Premium for 3 months \n\nRef#%s",
	6:  "Telegram Premium for 6 months \n\nRef#%s",
	12: "Telegram Premium for 1 year \n\nRef#%s",
}

// BuyTelegramPremium 支付成功后后台进行购买的流程
func buyTelegramPremium(tgUsername string, vipMonth int) (fragmentRefId string, err error) {
	result1, err := fragment.SearchPremiumGiftRecipient(tgUsername, vipMonth)
	if err != nil {
		return
	}

	result2, err := fragment.InitGiftPremium(result1.Found.Recipient, vipMonth)
	if err != nil {
		return
	}

	result3, err := fragment.GetGiftPremiumLink(result2.ReqId)
	if err != nil {
		return
	}

	info, err := fragment.GetTonPaymentInfo(result3.CheckParams.Id)
	if err != nil {
		return
	}

	receiverAddress := info.Body.Params.Messages[0].Address
	amount := info.Body.Params.Messages[0].Amount
	payload := info.Body.Params.Messages[0].Payload

	decodeBytes, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		return
	}
	arr := strings.Split(string(decodeBytes), "#")
	fragmentRef := arr[1]

	err = blockchain.Transfer(receiverAddress, amount, fmt.Sprintf(tonCommentFormats[vipMonth], fragmentRef))
	if err != nil {
		return
	}

	return
}

func updateSuccess(orderNo, fragmentRefId string) (err error) {
	err = query.Q.Transaction(func(tx *query.Query) error {
		var o, u = tx.Order, tx.User
		dbOrder, e := o.Where(o.OrderNo.Eq(orderNo)).First()
		if e != nil {
			return e
		}
		getUsdt := dbOrder.UsdtAmount - dbOrder.BaseAmount //代理赚的钱
		_, e = u.Where(u.ID.Eq(dbOrder.AgentUserID)).Update(u.Brokerage, u.Brokerage.Add(getUsdt))
		if e != nil {
			return e
		}
		_, e = o.Where(o.OrderNo.Eq(orderNo)).Updates(model.Order{
			FragmentRefID: &fragmentRefId,
			Status:        3,
			AgentStatus:   2,
		})
		if e != nil {
			return e
		}

		return e
	})

	return
}
