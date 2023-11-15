package svc

import (
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
