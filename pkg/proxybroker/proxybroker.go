package proxybroker

import (
	"container/heap"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/chneau/proxybroker/pkg/proxylist"
	"github.com/chneau/proxybroker/pkg/rate"
)

type ProxyBroker struct {
	// Proxies
	PriorityQueue PriorityQueue
	ProxyExist    map[string]bool
	mtxSelection  *sync.Mutex

	// Config
	SourceFn                   func() []string
	ProxyTesterFn              func(string) *Proxy
	LimitsPerDomain            map[string]*rate.Limit
	DurationBetweenSourceFetch time.Duration
	NumberOfParallelTest       int
}

func (pb *ProxyBroker) selection(req *http.Request) *Proxy {
	pb.mtxSelection.Lock()
	defer pb.mtxSelection.Unlock()
	bestWhen := time.Second
	bestProxy := (*Proxy)(nil)
	for _, proxy := range pb.PriorityQueue {
		if !proxy.IsReady {
			continue
		}
		when := proxy.LimitsPerDomain[req.Host].When()
		if when < bestWhen {
			bestWhen = when
			bestProxy = proxy
		}
	}
	time.Sleep(bestWhen)
	if bestProxy != nil {
		heap.Remove(&pb.PriorityQueue, bestProxy.index)
	}
	return bestProxy
}
func (pb *ProxyBroker) putBack(proxy *Proxy) {
	pb.mtxSelection.Lock()
	defer pb.mtxSelection.Unlock()
	pb.PriorityQueue.Push(proxy)
}

func (pb *ProxyBroker) Do(req *http.Request) (result []byte) {
	var proxy *Proxy
	for {
		proxy = pb.selection(req)
		if proxy == nil {
			continue
		}
		result = proxy.Do(req)
		log.Println("proxy", proxy.Name, proxy.index, "===>", string(result), len(result), result == nil)
		pb.putBack(proxy)
		if result == nil {
			continue
		}
		return result
	}
}

func (pb *ProxyBroker) WithDomainRateLimit(domain string, limit *rate.Limit) *ProxyBroker {
	pb.LimitsPerDomain[domain] = limit
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
		if _, exist := pb.ProxyExist[proxy]; !exist {
			client := pb.ProxyTesterFn(proxy)
			if client != nil {
				for k, v := range pb.LimitsPerDomain {
					client.LimitsPerDomain[k] = v.Clone()
				}
				pb.PriorityQueue.Push(client)
				pb.ProxyExist[proxy] = true
			}
		}
	}
}

func (pb *ProxyBroker) Init(waitN int) *ProxyBroker {
	newArrival := make(chan string)
	go pb.autoFetchSource(newArrival)
	for i := 0; i < pb.NumberOfParallelTest; i++ {
		go pb.tester(newArrival)
	}
	for {
		len := pb.PriorityQueue.Len()
		log.Println("Waiting ", len, "/", waitN)
		if len >= waitN {
			break
		}
		time.Sleep(time.Millisecond * 250)
	}
	return pb
}

func NewDefault() *ProxyBroker {
	pb := &ProxyBroker{
		SourceFn:                   proxylist.All,
		LimitsPerDomain:            map[string]*rate.Limit{},
		ProxyTesterFn:              ProxyTester,
		DurationBetweenSourceFetch: time.Minute,
		NumberOfParallelTest:       50,
		PriorityQueue:              PriorityQueue{},
		ProxyExist:                 map[string]bool{},
		mtxSelection:               &sync.Mutex{},
	}
	return pb
}
