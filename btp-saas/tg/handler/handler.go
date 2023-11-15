package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	tele "gopkg.in/telebot.v3"
)

const (
	RechargeKeyboardText    = "💰立即充值"
	MineKeyboardText        = "👤个人中心"
	AgentKeyboardText       = "💵代理赚钱"
	CooperationKeyboardText = "🤝商务合作"
)

const (
	CloseBtnId   = "CLOSE_BTN"
	SupportBtnId = "SUPPORT_BTN"

	ByMyselfBtnId  = "BY_MYSELF_BTN"
	GiftOtherBtnId = "GIFT_OTHER_BTN"

	BuyThreeMonthBtnId  = "BUY_THREE_MONTH_BTN"
	BuySixMonthBtnId    = "BUY_SIX_MONTH_BTN"
	BuyTwelveMonthBtnId = "BUY_TWELVE_MONTH_BTN"

	RechargeBtnId        = "RECHARGE_BTN"
	RechargeConfirmBtnId = "RECHARGE_INFO_CONFIRM_BTN"
)

var (
	RechargeKeyboard    = tele.Btn{Text: RechargeKeyboardText}
	MineKeyboard        = tele.Btn{Text: MineKeyboardText}
	AgentKeyboard       = tele.Btn{Text: AgentKeyboardText}
	CooperationKeyboard = tele.Btn{Text: CooperationKeyboardText}
)

var (
	CloseBtn    = tele.Btn{Unique: CloseBtnId, Text: "关闭"}
	SupportBtn  = tele.Btn{Unique: SupportBtnId, Text: "联系客服"}
	RechargeBtn = tele.Btn{Unique: RechargeBtnId, Text: "立即充值"}
)

func CloseHandler(ctx tele.Context) error {
	return ctx.Delete()
}

func SupportHandler(ctx tele.Context) error {
	var u = query.User
	agentUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("[db] query fail. %v\n", err)
		return err
	}
	replyFormat := "t.me/%s"
	reply := EscapeText(tele.ModeMarkdownV2, fmt.Sprintf(replyFormat, *agentUser.TgUsername))
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
