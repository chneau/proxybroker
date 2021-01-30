package proxybroker

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var ipRegex = regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)

func ProxyTester(proxy string) bool {
	req, _ := http.NewRequest("GET", "http://api.ipify.org/", nil)
	proxyUrl, _ := url.Parse(proxy)
	client := &http.Client{
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	ip, _ := ioutil.ReadAll(resp.Body)
	success := ipRegex.Match(ip)
	return success
}
