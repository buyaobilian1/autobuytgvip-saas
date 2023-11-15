package global

import (
	"net/http"
	"net/url"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/config"
)

var (
	Conf config.Config
)

func NewProxyHttpClient(proxyUrl string) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}
	httpTransport := &http.Transport{
		Proxy: proxy,
	}
	return &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * 10,
	}
}
