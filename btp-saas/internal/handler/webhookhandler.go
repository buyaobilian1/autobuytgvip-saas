package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/tg"
	"github.com/zeromicro/go-zero/rest/httpx"
	tele "gopkg.in/telebot.v3"
)

func WebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		botIdStr := strings.TrimPrefix(r.URL.Path, "/webhook/")
		botId, _ := strconv.ParseInt(botIdStr, 10, 64)
		bot, err := tg.GetBotInstance(botId)
		if err != nil {
			log.Printf("[bot %d] get bot instance fail. %v\n", botId, err)
			httpx.Ok(w)
			return
		}

		var update tele.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Println(fmt.Errorf("cannot decode update: %v", err))
			httpx.Ok(w)
			return
		}
		bot.ProcessUpdate(update)
		httpx.Ok(w)
	}
}
