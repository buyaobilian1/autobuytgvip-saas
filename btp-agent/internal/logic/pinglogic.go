package logic

import (
	"context"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping(req *types.PingRequest) (resp *types.PingResponse, err error) {

	return
}
