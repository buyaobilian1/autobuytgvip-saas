package fragment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type SearchPremiumGiftRecipientRequest struct {
	Cookie string `header:"cookie"`
	Query  string `form:"query"`
	Months int    `form:"months"`
	Method string `form:"method"`
}

type SearchPremiumGiftRecipientResponse struct {
	Ok    bool                                   `json:"ok"`
	Error string                                 `json:"error"`
	Found SearchPremiumGiftRecipientResponseBody `json:"found"`
}

type SearchPremiumGiftRecipientResponseBody struct {
	Myself    bool   `json:"myself"`
	Recipient string `json:"recipient"`
	Photo     string `json:"photo"`
	Name      string `json:"name"`
}

// SearchPremiumGiftRecipient telegram用户名查询
func SearchPremiumGiftRecipient(username string, duration int) (result SearchPremiumGiftRecipientResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := SearchPremiumGiftRecipientRequest{
		Cookie: global.Conf.AppConf.Cookie,
		Query:  username,
		Months: duration,
		Method: "searchPremiumGiftRecipient",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = SearchPremiumGiftRecipientResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil

}

type InitGiftPremiumRequest struct {
	Cookie    string `header:"cookie"`
	Recipient string `form:"recipient"`
	Months    int    `form:"months"`
	Method    string `form:"method"`
}

type InitGiftPremiumResponse struct {
	ReqId     string `json:"req_id"`
	Myself    bool   `json:"myself"`
	Amount    string `json:"amount"`
	ItemTitle string `json:"item_title"`
	Content   string `json:"content"`
	Button    string `json:"button"`
}

func InitGiftPremium(recipient string, duration int) (result InitGiftPremiumResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := InitGiftPremiumRequest{
		Cookie:    global.Conf.AppConf.Cookie,
		Recipient: recipient,
		Months:    duration,
		Method:    "initGiftPremiumRequest",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = InitGiftPremiumResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type GetGiftPremiumLinkRequest struct {
	Cookie     string `header:"cookie"`
	Id         string `form:"id"`
	ShowSender int    `form:"show_sender"`
	Months     int    `form:"months"`
	Method     string `form:"method"`
}

type GetGiftPremiumLinkResponse struct {
	Ok          bool        `json:"ok"`
	Link        string      `json:"link"`
	QrLink      string      `json:"qr_link"`
	CheckMethod string      `json:"check_method"`
	CheckParams CheckParams `json:"check_params"`
	ExpireAfter int         `json:"expire_after"`
}
type CheckParams struct {
	Id string `json:"id"`
}

func GetGiftPremiumLink(reqId string) (result GetGiftPremiumLinkResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := GetGiftPremiumLinkRequest{
		Cookie:     global.Conf.AppConf.Cookie,
		Id:         reqId,
		ShowSender: 0,
		Method:     "getGiftPremiumLink",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = GetGiftPremiumLinkResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type GetTonPaymentInfoRequest struct {
	Cookie string `header:"cookie"`
}

type GetTonPaymentInfoResponse struct {
	Version string `json:"version"`
	Body    struct {
		Type   string `json:"type"`
		Params struct {
			ValidUntil int `json:"valid_until"`
			Messages   []struct {
				Address string `json:"address"`
				Amount  uint64 `json:"amount"`
				Payload string `json:"payload"`
			} `json:"messages"`
			Source string `json:"source"`
		} `json:"params"`
		ResponseOptions struct {
			CallbackURL string `json:"callback_url"`
			Broadcast   bool   `json:"broadcast"`
		} `json:"response_options"`
		ExpiresSec int `json:"expires_sec"`
	} `json:"body"`
}

// GetTonPaymentInfo GET获取收款地址和 payload参数
func GetTonPaymentInfo(id string) (result GetTonPaymentInfoResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/tonkeeper/rawRequest?id=%s&qr=1", id)
	req := GetTonPaymentInfoRequest{
		Cookie: global.Conf.AppConf.Cookie,
	}
	resp, err := httpc.Do(context.Background(), http.MethodGet, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = GetTonPaymentInfoResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type CheckOrderRequest struct {
	Cookie string `header:"cookie"`
	Id     string `form:"id"`
	Method string `form:"method"`
}

type CheckOrderResponse struct {
	Confirmed bool `json:"confirmed"`
}

// CheckOrder 检查订单是否成功
func CheckOrder(id string) (result CheckOrderResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := CheckOrderRequest{
		Cookie: global.Conf.AppConf.Cookie,
		Id:     id,
		Method: "checkReq",
	}
	resp, err := httpc.Do(context.Background(), http.MethodGet, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = CheckOrderResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}
