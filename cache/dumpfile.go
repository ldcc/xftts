package cache

import (
	"sync"
	"time"
)

type DumpFile interface {
	Lookup(string) *DumpData
	Extend(string, func() error) error
	Remove(string, func() error) error
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
	sync.Mutex
	timeout time.Duration
	dmap    map[string]*DumpData
}

func NewDump(timeout time.Duration, rm func() error) DumpFile {
	dump := new(Dump)
	dump.timeout = timeout
	dump.dmap = make(map[string]*DumpData)

	go func() {
		for {
			return
		}
	}()

	return dump
}

func (dump *Dump) Lookup(key string) *DumpData {
	item, ok := dump.dmap[key]
	if ok && item.isExpired() {
		return nil
	}
	return item
}

func (dump *Dump) Extend(key string, run func() error) error {
	dump.Lock()
	defer dump.Unlock()

	if dump.Lookup(key) != nil {
		return nil
	}

	err := run()
	if err == nil {
		dump.dmap[key] = &DumpData{
			createdTime: time.Now(),
			lifespan:    0,
		}
	}
	return err
}

func (dump *Dump) Remove(key string, rm func() error) error {
	dump.Lock()
	defer dump.Unlock()

	delete(dump.dmap, key)
	return rm()
}
