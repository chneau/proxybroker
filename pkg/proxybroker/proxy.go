package proxybroker

import (
	"github.com/chneau/proxybroker/pkg/ratelimit"
	"net/http"
)

type Proxy struct {
	*http.Client
	Score           int
	LimitsPerDomain map[string]ratelimit.Limit
}
