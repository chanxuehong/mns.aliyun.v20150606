package internal

import (
	"net/url"
	"sync"
	"sync/atomic"
	"unsafe"
)

var ParseURL = (&urlCache{}).Parse

type urlCache struct {
	sync.Mutex
	mapPtr unsafe.Pointer // *urlMap
}

type urlMap map[string]*url.URL // map[rawurl]*url.URL

func (cache *urlCache) Parse(rawurl string) (*url.URL, error) {
	var m urlMap

	if p := (*urlMap)(atomic.LoadPointer(&cache.mapPtr)); p != nil {
		m = *p
		if u := m[rawurl]; u != nil {
			return u, nil
		}
	}

	cache.Lock()
	defer cache.Unlock()

	if p := (*urlMap)(atomic.LoadPointer(&cache.mapPtr)); p != nil {
		m = *p
		if u := m[rawurl]; u != nil {
			return u, nil
		}
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	m2 := make(urlMap, len(m)+1)
	for k, v := range m {
		m2[k] = v
	}
	m2[rawurl] = u

	atomic.StorePointer(&cache.mapPtr, unsafe.Pointer(&m2))
	return u, nil
}
