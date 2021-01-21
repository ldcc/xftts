package cache

import (
	"sync"
	"time"
)

type DumpFile interface {
	Lookup(string) *DumpData
	Extend(string, func() error) error
}

type DumpData struct {
	createdTime time.Time
	lifespan    time.Duration
}

func (dd *DumpData) isExpired() bool {
	// 0 为永不过期
	if dd.lifespan == 0 {
		return false
	}
	return time.Now().Sub(dd.createdTime) > dd.lifespan
}

type Dump struct {
	sync.RWMutex
	timeout time.Duration
	items   map[string]*DumpData
	drop    func(string)
}

/**
 * timeout 超时时间
 */
func NewDump(timeout time.Duration, drop func(string)) DumpFile {
	dump := new(Dump)
	dump.timeout = timeout
	dump.items = make(map[string]*DumpData)
	dump.drop = drop

	go func() {
		for {
			<-time.After(timeout)
			keys := dump.expiredKeys()
			if len(keys) != 0 {
				dump.clearItems(keys)
			}
		}
	}()

	return dump
}

func (dump *Dump) Lookup(key string) *DumpData {
	item, ok := dump.items[key]
	if ok && item.isExpired() {
		return nil
	}
	return item
}

func (dump *Dump) Extend(key string, dumping func() error) error {
	dump.Lock()
	defer dump.Unlock()

	if dump.Lookup(key) != nil {
		return nil
	}

	err := dumping()
	if err == nil {
		dump.items[key] = &DumpData{
			createdTime: time.Now(),
			lifespan:    dump.timeout,
		}
	}
	return err
}

func (dump *Dump) clearItems(keys []string) {
	dump.Lock()
	defer dump.Unlock()

	for _, key := range keys {
		delete(dump.items, key)
		dump.drop(key)
	}
}

func (dump *Dump) expiredKeys() (keys []string) {
	dump.RLock()
	defer dump.RUnlock()
	for key, itm := range dump.items {
		if itm.isExpired() {
			keys = append(keys, key)
		}
	}
	return
}
