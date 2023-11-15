package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/config"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq/handle"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/pkg/epusdt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type CreateEpusdtPaymentRequest struct {
	OrderId     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	RedirectUrl string  `json:"redirect_url"`
	NotifyUrl   string  `json:"notify_url"`
	Signature   string  `json:"signature"`
}

type CreateEpusdtPaymentResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	RequestId  string `json:"request_id"`
	Data       struct {
		TradeId        string  `json:"trade_id"`
		OrderId        string  `json:"order_id"`
		Amount         float64 `json:"amount"`
		ActualAmount   float64 `json:"actual_amount"`
		Token          string  `json:"token"`
		ExpirationTime int64   `json:"expiration_time"`
		PaymentUrl     string  `json:"payment_url"`
	} `json:"data"`
}

func CreatePayOrder(conf config.PayConf, order *model.Order) (result CreateEpusdtPaymentResponse, err error) {
	api := fmt.Sprintf("%s/api/v1/order/create-transaction", conf.BaseApi)
	// 这里CNY与USDT的价格先写死，后面做动态的
	// todo 生产环境里强制设置了汇率是1，所以这里不需要考虑汇率了
	amount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", order.UsdtAmount), 64)

	req := CreateEpusdtPaymentRequest{
		OrderId:   order.OrderNo,
		Amount:    amount,
		NotifyUrl: conf.NotifyUrl,
	}
	var m = make(map[string]interface{})
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &m)
	signature, _ := epusdt.Sign(m, conf.ApiToken)
	req.Signature = signature

	resp, err := httpc.Do(context.Background(), http.MethodPost, api, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = CreateEpusdtPaymentResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	if result.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("epusdt order create fail, code = %d", result.StatusCode))
		return
	}

	// 本地订单创建
	order.UsdtAmount = result.Data.ActualAmount
	var orderDao = query.Order
	err = orderDao.Create(order)
	if err != nil {
		return
	}
	task, _ := handle.NewOrderExpirationTask(order.OrderNo)
	duration := time.Unix(result.Data.ExpirationTime, 0).Sub(time.Now())
	fmt.Println(duration)
	// 比支付平台早10秒钟关闭订单
	_, _ = mq.QueueClient.Enqueue(task, asynq.ProcessIn(duration-10*time.Second))

	return result, nil
}

func CreateEpusdtPayment(orderNo string, usdtAmount float64, notifyUrl string) (result CreateEpusdtPaymentResponse, err error) {
	var conf = global.Conf.PayConf
	api := fmt.Sprintf("%s/api/v1/order/create-transaction", conf.BaseApi)
	// 这里CNY与USDT的价格先写死，后面做动态的
	// todo 生产环境里强制设置了汇率是1，所以这里不需要考虑汇率了
	amount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", usdtAmount), 64)

	req := CreateEpusdtPaymentRequest{
		OrderId:   orderNo,
		Amount:    amount,
		NotifyUrl: notifyUrl,
	}
	var m = make(map[string]interface{})
	b, _ := json.Marshal(req)
	_ = json.Unmarshal(b, &m)
	signature, _ := epusdt.Sign(m, conf.ApiToken)
	req.Signature = signature

	resp, err := httpc.Do(context.Background(), http.MethodPost, api, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = CreateEpusdtPaymentResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	if result.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("epusdt order create fail, code = %d", result.StatusCode))
		return
	}

	return
}
