package handler

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
)

const (
	OrderKeyboardText   = "💎订单"
	PriceKeyboardText   = "💵价格设置"
	FinanceKeyboardText = "💰财务"
	BotKeyboardText     = "⚙️设置"
)

const (
	CloseBtnId    = "CLOSE_BTN"
	SupportBtnId  = "SUPPORT_BTN"
	WithDrawBtnId = "WITHDRAW_BTN"
	RechargeBtnId = "RECHARGE_BTN"
)

const (
	BotTokenSettingBtnId   = "BOT_TOKEN_SETTING_BTN"
	AgentPriceSettingBtnId = "AGENT_PRICE_SETTING_BTN"
	OrderPagePrevBtnId     = "ORDER_PAGE_PREV_BTN"
	OrderPageNextBtnId     = "ORDER_PAGE_NEXT_BTN"

	OrderDetailBtn     = "ORDER_DETAIL_BTN"
	OrderDetailBackBtn = "ORDER_DETAIL_BACK_BTN"
)

var (
	FinanceKeyboard = tele.Btn{Text: FinanceKeyboardText}
	OrderKeyboard   = tele.Btn{Text: OrderKeyboardText}
	BotKeyboard     = tele.Btn{Text: BotKeyboardText}
	PriceKeyboard   = tele.Btn{Text: PriceKeyboardText}
)

var (
	CloseBtn    = tele.Btn{Unique: CloseBtnId, Text: "关闭"}
	SupportBtn  = tele.Btn{Unique: SupportBtnId, Text: "联系客服"}
	WithDrawBtn = tele.Btn{Unique: WithDrawBtnId, Text: "提现"}

	RechargeBtn = tele.Btn{Unique: RechargeBtnId, Text: "立即充值"}

	BotTokenSettingBtn   = tele.Btn{Unique: BotTokenSettingBtnId, Text: "🤖机器人Token设置"}
	AgentPriceSettingBtn = tele.Btn{Unique: AgentPriceSettingBtnId, Text: "💵代理销售价格设置"}
)

func CloseHandler(ctx tele.Context) error {
	return ctx.Delete()
}

func SupportHandler(ctx tele.Context) error {
	return ctx.Send(EscapeText(tele.ModeMarkdownV2, "t.me/feijige120"))
}

func HelpHandler(ctx tele.Context) error {
	reply := `
/start 开始使用
/help 帮助
/token 设置机器人Token
/price3 设置3个月会员销售价
/price6 设置6个月会员销售价
/price12 设置12个月会员销售价
/orders 代理订单查询
/address 设置提现钱包地址
/withdraw 发起提现
`

	return ctx.Send(reply)
}

// EscapeText 机器人文本处理
func EscapeText(parseMode string, text string) string {
	var replacer *strings.Replacer

	if parseMode == tele.ModeHTML {
		replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	} else if parseMode == tele.ModeMarkdown {
		replacer = strings.NewReplacer("_", "\\_", "*", "\\*", "`", "\\`", "[", "\\[")
	} else if parseMode == tele.ModeMarkdownV2 {
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	} else {
		return ""
	}

	return replacer.Replace(text)
}

func Float64Format(money float64) string {
	moneyStr := fmt.Sprintf("%f", money)
	for strings.HasSuffix(moneyStr, "0") {
		moneyStr = strings.TrimSuffix(moneyStr, "0")
	}
	if strings.HasSuffix(moneyStr, ".") {
		moneyStr = strings.TrimSuffix(moneyStr, ".")
	}
	//moneyStr = strings.ReplaceAll(moneyStr, ".", "\\.")

	return moneyStr
}
