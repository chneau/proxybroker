package proxybroker

import (
	"log"
	"net/http"
	"time"

	"github.com/chneau/proxybroker/pkg/proxylist"
	"github.com/chneau/proxybroker/pkg/ratelimit"
)

type ProxyBroker struct {
	// Config
	SourceFn                   func() []string
	ProxyTesterFn              func(string) bool
	LimitsPerDomain            map[string]ratelimit.Limit
	DurationBetweenSourceFetch time.Duration
	NumberOfParallelTest       int
}

func (pb *ProxyBroker) Do(req *http.Request) []byte {
	return nil
}

func (pb *ProxyBroker) WithDomainRateLimit(domain string, limit *ratelimit.Limit) *ProxyBroker {
	pb.LimitsPerDomain[domain] = *limit
	return pb
}

func (pb *ProxyBroker) WithSourceFn(sourceFn func() []string) *ProxyBroker {
	pb.SourceFn = sourceFn
	return pb
}

func (pb *ProxyBroker) autoFetchSource(newArrival chan string) {
	for {
		proxies := pb.SourceFn()
		for _, proxy := range proxies {
			newArrival <- proxy
		}
		time.Sleep(pb.DurationBetweenSourceFetch)
	}
}

func (pb *ProxyBroker) tester(newArrival chan string) {
	for proxy := range newArrival {
		if pb.ProxyTesterFn(proxy) {
			log.Println(proxy)
		} else {
			log.Println("failed")
		}
	}
}

func (pb *ProxyBroker) Init() *ProxyBroker {
	newArrival := make(chan string)
	go pb.autoFetchSource(newArrival)
	for i := 0; i < pb.NumberOfParallelTest; i++ {
		go pb.tester(newArrival)
	}
	return pb
}

func NewDefault() *ProxyBroker {
	pb := &ProxyBroker{
		SourceFn:                   proxylist.All,
		LimitsPerDomain:            map[string]ratelimit.Limit{},
		ProxyTesterFn:              ProxyTester,
		DurationBetweenSourceFetch: time.Minute,
		NumberOfParallelTest:       10,
	}
	return pb
}
