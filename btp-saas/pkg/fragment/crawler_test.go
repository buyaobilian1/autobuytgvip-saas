package fragment

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestSearchPremiumGiftRecipient(t *testing.T) {
	duration := 3

	result1, _ := SearchPremiumGiftRecipient("@king3670", duration)
	fmt.Printf("查询Telegram用户信息：%+v\n", result1)

	result2, _ := InitGiftPremium(result1.Found.Recipient, duration)
	fmt.Printf("初始化赠送Telegram会员请求：%+v\n", result2)

	result3, _ := GetGiftPremiumLink(result2.ReqId)
	fmt.Printf("获取Telegr会员：%+v\n", result3)

	info, _ := GetTonPaymentInfo(result3.CheckParams.Id)
	fmt.Printf("支付信息：%+v\n", info)

	receiverAddress := info.Body.Params.Messages[0].Address
	amount := info.Body.Params.Messages[0].Amount
	payload := info.Body.Params.Messages[0].Payload

	decodeBytes, _ := base64.RawStdEncoding.DecodeString(payload)
	arr := strings.Split(string(decodeBytes), "#")
	orderSN := arr[1]

	commentFormatter := `Telegram Premium for 3 months 

Ref#%s`
	fmt.Printf("address: %s\namount: %d\ncomment:%s\n", receiverAddress, amount, fmt.Sprintf(commentFormatter, orderSN))
}
