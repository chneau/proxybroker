package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chneau/proxybroker/pkg/proxybroker"
	"github.com/chneau/proxybroker/pkg/ratelimit"
	"github.com/elazarl/goproxy"
)

func main() {
	go SetUpLocalProxy()
	pb := proxybroker.
		NewDefault().
		WithDomainRateLimit("ip-api.com", ratelimit.New().WithLimit(1, time.Second)).
		WithSourceFn(func() []string { return []string{"127.0.0.1:5000"} }).
		Init()
	_ = pb
	time.Sleep(time.Hour)
}

func SetUpLocalProxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Println("Set up local proxy")
	log.Fatal(http.ListenAndServe(":5000", proxy))
}
