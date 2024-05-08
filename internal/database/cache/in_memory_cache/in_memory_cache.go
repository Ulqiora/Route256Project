package in_memory_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var defaultTimeout = 10 * time.Minute

type InMemoryCache struct {
	data map[string]struct {
		info []byte
		ttl  time.Time
	}
	mtx sync.RWMutex
}

func (i *InMemoryCache) Ping(_ context.Context) error {
	return nil
}

func (i *InMemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	data, ok := i.data[key]
	if !ok {
		return nil, fmt.Errorf("object not exists with key : %s", key)
	}
	if time.Now().After(data.ttl) {
		return nil, errors.New("time out")
	}
	return data.info, nil
}

func (i *InMemoryCache) Set(_ context.Context, key string, a any) error {
	marshal, err := json.Marshal(a)
	if err != nil {
		return err
	}
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.data[key] = struct {
		info []byte
		ttl  time.Time
	}{info: marshal, ttl: time.Now().Add(defaultTimeout)}
	return nil
}

func (i *InMemoryCache) SetWithTTL(_ context.Context, ttl time.Duration, key string, a any) error {
	marshal, err := json.Marshal(a)
	if err != nil {
		return err
	}
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.data[key] = struct {
		info []byte
		ttl  time.Time
	}{info: marshal, ttl: time.Now().Add(ttl)}
	return nil
}

func (i *InMemoryCache) SetMulti(_ context.Context, a map[string]any) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	for key, value := range a {
		marshal, err := json.Marshal(value)
		if err != nil {
			return err
		}
		i.data[key] = struct {
			info []byte
			ttl  time.Time
		}{info: marshal, ttl: time.Now().Add(defaultTimeout)}
	}
	return nil
}

func (i *InMemoryCache) SetMultiWithTTL(_ context.Context, ttl time.Duration, a map[string]any) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	for key, value := range a {
		marshal, err := json.Marshal(value)
		if err != nil {
			return err
		}
		i.data[key] = struct {
			info []byte
			ttl  time.Time
		}{info: marshal, ttl: time.Now().Add(ttl)}
	}
	return nil
}

func (i *InMemoryCache) Run() {
	go func() {
		interval := defaultTimeout / 2
		t := time.NewTicker(interval)
		for {
			select {
			case <-t.C:
				now := time.Now()
				for key, v := range i.data {
					sub := v.ttl.Sub(now).Hours()
					if sub > 12 {
						delete(i.data, key)
					}
				}
			default:
			}
		}
	}()

}

func New() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]struct {
			info []byte
			ttl  time.Time
		}, 100),
		mtx: sync.RWMutex{},
	}
}
