package main

import (
	"log"
	"time"

	"github.com/chneau/proxybroker/pkg/proxybroker"
	"github.com/chneau/proxybroker/pkg/ratelimit"
)

func main() {
	log.Println("Helelo world !")
	pb := proxybroker.
		NewDefault().
		WithDomainRateLimit("ip-api.com", ratelimit.New().WithLimit(1, time.Second))
	_ = pb
}
