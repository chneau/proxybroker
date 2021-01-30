package proxybroker

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/chneau/proxybroker/pkg/rate"
)

type Proxy struct {
	Client          *http.Client
	Times           []time.Duration
	LimitsPerDomain map[string]rate.Limit
	mtx             *sync.Mutex
}

// sync allows to write `defer proxy.sync()()` one liner
func (proxy *Proxy) sync() func() {
	proxy.mtx.Lock()
	return proxy.mtx.Unlock
}

// MeanTime returns the mean time it took the last 10 requests
func (proxy *Proxy) MeanTime() time.Duration {
	defer proxy.sync()()
	mean := time.Duration(0)
	for _, time := range proxy.Times {
		mean += time
	}
	mean /= time.Duration(len(proxy.Times))
	return mean
}

func (proxy *Proxy) Do(req *http.Request) []byte {
	defer proxy.sync()()
	start := time.Now()
	resp, err := proxy.Client.Do(req)
	if err != nil {
		proxy.Times = append(proxy.Times, time.Second*30)
		return nil
	}
	defer resp.Body.Close()
	bb, _ := ioutil.ReadAll(resp.Body)
	proxy.Times = append(proxy.Times, time.Since(start))
	if len(proxy.Times) > 10 {
		proxy.Times = proxy.Times[len(proxy.Times)-10:]
	}
	return bb
}

func NewProxy(proxy string) *Proxy {
	proxyUrl, _ := url.Parse(proxy)
	p := &Proxy{
		Client: &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				DisableKeepAlives: true,
				Proxy:             http.ProxyURL(proxyUrl),
			},
		},
		mtx: &sync.Mutex{},
	}
	return p
}
