package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chneau/proxybroker/pkg/proxybroker"
	"github.com/chneau/proxybroker/pkg/rate"
	"github.com/elazarl/goproxy"
)

func main() {
	go SetUpLocalProxy()
	pb := proxybroker.
		NewDefault().
		WithDomainRateLimit("api.ipify.org", rate.NewLimit().WithConstraint(1, time.Second*3)).
		Init(5)
	_ = pb
	for {
		go func() {
			req, _ := http.NewRequest("GET", "http://api.ipify.org/", nil)
			ip := pb.Do(req)
			log.Println(string(ip))
		}()
		time.Sleep(time.Millisecond * 250)
	}
}

func SetUpLocalProxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Println("Set up local proxy")
	log.Fatal(http.ListenAndServe(":5000", proxy))
}
