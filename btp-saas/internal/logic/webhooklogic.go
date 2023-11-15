package logic

import (
	"context"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type WebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookLogic {
	return &WebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookLogic) Webhook(req *types.WebhookRequest) error {
	// todo: add your logic here and delete this line

	return nil
}
