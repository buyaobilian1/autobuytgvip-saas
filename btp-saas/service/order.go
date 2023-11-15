package service

import (
	"fmt"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq/handle"
	"github.com/hibiken/asynq"
)

type CreateOrderResponse struct {
	Token        string  `json:"token"`
	ActualAmount float64 `json:"actual_amount"`
}

func CreateOrder(order *model.Order) (result CreateOrderResponse, err error) {
	var o = query.Order
	notifyUrl := fmt.Sprintf(global.Conf.PayConf.NotifyUrl, "order")
	payment, err := CreateEpusdtPayment(order.OrderNo, order.UsdtAmount, notifyUrl)
	if err != nil {
		return
	}

	order.UsdtAmount = payment.Data.ActualAmount
	err = o.Create(order)
	if err != nil {
		return
	}
	task, _ := handle.NewOrderExpirationTask(order.OrderNo)
	_, _ = mq.QueueClient.Enqueue(task, asynq.ProcessIn(time.Minute*time.Duration(global.Conf.AppConf.OrderExpireMinute)))

	return CreateOrderResponse{
		Token:        payment.Data.Token,
		ActualAmount: payment.Data.ActualAmount,
	}, nil
}
