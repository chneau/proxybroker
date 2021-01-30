package proxylist

import (
	"sync"
	"testing"
)

func TestAll(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(len(ProxyLists))
	for name := range ProxyLists {
		name := name
		go func() {
			defer wg.Done()
			prx := ProxyLists[name]
			found := prx()
			if len(found) == 0 {
				t.Error(name, len(found))
			}
		}()
	}
	wg.Wait()
}
