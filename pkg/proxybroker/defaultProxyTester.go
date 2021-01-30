package proxybroker

import (
	"net/http"
	"regexp"
)

var ipRegex = regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)

func ProxyTester(proxy string) *Proxy {
	req, _ := http.NewRequest("GET", "http://api.ipify.org/", nil)
	client := NewProxy(proxy)
	ip := client.Do(req)
	success := ipRegex.Match(ip)
	if !success {
		return nil
	}
	return client
}
