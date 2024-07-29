package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func removeOtherKeys(keysToKeep map[interface{}]bool, m *sync.Map) {
	var allKeys []interface{}
	m.Range(func(k interface{}, v interface{}) bool {
		allKeys = append(allKeys, k)
		return true
	})
	for _, k := range allKeys {
		if !keysToKeep[k] {
			m.Delete(k)
		}
	}
}

func deletePortInfo(port uint16) {
	pid, _ := PortConfMap.Load(port)
	hashMap.Store(pid.(string), false)
	PortPidMap.Delete(port)
	PortConfMap.Delete(port)
	freePorts.Insert(port)
	threadsNow.Add(-1)
}

func httpGET(url string, maxAttempts int) string {
	var attempt = 0
	var resp *http.Response
	var err error
	for attempt < maxAttempts {
		attempt++
		resp, err = http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(20 * time.Second)
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return string(b)
		}
	}
	return ""
}
