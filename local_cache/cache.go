package local_cache

import (
	"fmt"
	"sync"
	"time"
)

type SafeMap struct {
	mu sync.Mutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, ok := sm.m[key]
	return value, ok
}

func S() {
	sm := NewSafeMap()

	// 更新Map的函数
	updateMap := func() {
		sm.Set("key", time.Now().Format(time.RFC3339))
	}

	// 每分钟更新一次Map
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateMap()
			value, ok := sm.Get("key")
			if ok {
				fmt.Println("Updated value:", value)
			}
		}
	}
}
