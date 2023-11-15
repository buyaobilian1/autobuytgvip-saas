package main

import (
	"flag"
	"fmt"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/config"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/handler"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/svc"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/tg"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/btp-agent.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	logx.Disable()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	global.Conf = c

	dao.Start(c)
	tg.Start(c.AppConf)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
