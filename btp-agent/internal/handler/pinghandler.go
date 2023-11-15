package handler

import (
	"net/http"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/logic"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPingLogic(r.Context(), svcCtx)
		resp, err := l.Ping(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
