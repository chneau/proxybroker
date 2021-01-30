package proxybroker

import (
	"net/http"

	"github.com/chneau/proxybroker/pkg/proxylist"
	"github.com/chneau/proxybroker/pkg/ratelimit"
)

type ProxyBroker struct {
	ProxiesGiver    func() []string
	LimitsPerDomain map[string]ratelimit.Limit
}

func (pb *ProxyBroker) Do(req *http.Request) []byte {
	return nil
}

func (pb *ProxyBroker) WithDomainRateLimit(domain string, limit *ratelimit.Limit) *ProxyBroker {
	pb.LimitsPerDomain[domain] = *limit
	return pb
}

func NewDefault() *ProxyBroker {
	pb := &ProxyBroker{
		ProxiesGiver:    proxylist.All,
		LimitsPerDomain: map[string]ratelimit.Limit{},
	}
	return pb
}
