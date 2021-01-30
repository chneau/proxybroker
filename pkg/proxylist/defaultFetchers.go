package proxylist

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	ProxyLists = map[string]func() []string{
		"ProxiesFromClarketm":       ProxiesFromClarketm,
		"ProxiesFromDailyFreeProxy": ProxiesFromDailyFreeProxy,
		"ProxiesFromDailyProxy":     ProxiesFromDailyProxy,
		"ProxiesFromFate0":          ProxiesFromFate0,
		"ProxiesFromSmallSeoTools":  ProxiesFromSmallSeoTools,
		"ProxiesFromSunny9577":      ProxiesFromSunny9577,
		"ProxiesFromTheSpeedX":      ProxiesFromTheSpeedX,
	}
	ipPortRegex = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+:\d+`)
)

// ProxiesFromDailyFreeProxy returns proxies from https://www.dailyfreeproxy.com/.
func ProxiesFromDailyFreeProxy() []string {
	resp, err := http.Get("https://www.dailyfreeproxy.com/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}
	urls := []string{}
	wg := sync.WaitGroup{}
	doc.Find("h3 > a").Each(func(i int, s *goquery.Selection) {
		next := s.AttrOr("href", "")
		if next == "" {
			return
		}
		if !strings.Contains(strings.ToLower(s.Text()), "http") {
			return
		}
		wg.Add(1)
		defer wg.Done()
		resp, err := http.Get(next)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		bb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		str := string(bb)
		proxies := ipPortRegex.FindAllString(str, -1)
		for _, p := range proxies {
			urls = append(urls, "http://"+p)
		}
	})
	wg.Wait()
	return urls
}

// ProxiesFromSmallSeoTools returns proxies from https://smallseotools.com/free-proxy-list/.
func ProxiesFromSmallSeoTools() []string {
	resp, err := http.Get("https://smallseotools.com/free-proxy-list/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	str := string(bb)
	proxies := ipPortRegex.FindAllString(str, -1)
	urls := make([]string, len(proxies))
	for i := range proxies {
		urls[i] = "http://" + proxies[i]
	}
	return urls
}

// ProxiesFromDailyProxy returns proxies from https://proxy-daily.com/.
func ProxiesFromDailyProxy() []string {
	resp, err := http.Get("https://proxy-daily.com/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}
	httpList := doc.Find(".centeredProxyList.freeProxyStyle").First().Text()
	proxies := ipPortRegex.FindAllString(httpList, -1)
	urls := make([]string, len(proxies))
	for i := range proxies {
		urls[i] = "http://" + proxies[i]
	}
	return urls
}

// ProxiesFromClarketm returns proxies from clarketm/proxy-list.
func ProxiesFromClarketm() []string {
	resp, err := http.Get("https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	str := string(bb)
	proxies := ipPortRegex.FindAllString(str, -1)
	urls := make([]string, len(proxies))
	for i := range proxies {
		urls[i] = "http://" + proxies[i]
	}
	return urls
}

// ProxiesFromTheSpeedX returns proxies from hookzof/socks5_list.
func ProxiesFromTheSpeedX() []string {
	resp, err := http.Get("https://github.com/TheSpeedX/PROXY-List/blob/master/http.txt")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	str := string(bb)
	proxies := ipPortRegex.FindAllString(str, -1)
	urls := make([]string, len(proxies))
	for i := range proxies {
		urls[i] = "http://" + proxies[i]
	}
	return urls
}

type proxy struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// ProxiesFromSunny9577 returns proxies from sunny9577/proxy-scraper.
func ProxiesFromSunny9577() []string {
	resp, err := http.Get("https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/proxies.json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	proxies := &struct {
		LastUpdated string  `json:"lastUpdated"`
		Proxynova   []proxy `json:"proxynova"`
		Usproxy     []proxy `json:"usproxy"`
		Hidemyname  []proxy `json:"hidemyname"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(proxies)
	if err != nil {
		return nil
	}
	urls := []string{}
	for _, proxy := range proxies.Proxynova {
		urls = append(urls, "http://"+proxy.IP+":"+proxy.Port)
	}
	for _, proxy := range proxies.Usproxy {
		urls = append(urls, "http://"+proxy.IP+":"+proxy.Port)
	}
	for _, proxy := range proxies.Hidemyname {
		urls = append(urls, "http://"+proxy.IP+":"+proxy.Port)
	}
	return urls
}

// ProxiesFromFate0 returns proxies from fate0/proxylist.
func ProxiesFromFate0() []string {
	resp, err := http.Get("https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	proxy := &struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{}
	urls := []string{}
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), proxy)
		if err != nil {
			return nil
		}
		urls = append(urls, "http://"+proxy.Host+":"+strconv.Itoa(proxy.Port))
	}
	return urls
}

func All() []string {
	allProxies := map[string]struct{}{}
	wg := sync.WaitGroup{}
	wg.Add(len(ProxyLists))
	for name := range ProxyLists {
		name := name
		go func() {
			defer wg.Done()
			proxies := ProxyLists[name]()
			for _, proxy := range proxies {
				allProxies[proxy] = struct{}{}
			}
		}()
	}
	wg.Wait()
	proxies := []string{}
	for proxy := range allProxies {
		proxies = append(proxies, proxy)
	}
	sort.Strings(proxies)
	return proxies
}

// TODO: check if other exist here https://github.com/topics/proxy-list?o=desc&s=updated
