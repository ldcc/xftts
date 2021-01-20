package cache

import "sync"

type DumpFile interface {
	Lookup(string) interface{}
	Extend(string)
}

type Dump struct {
	sync.Mutex
	dmap map[string]interface{}
}

func NewDump() DumpFile {
	dump := new(Dump)
	dump.dmap = make(map[string]interface{})
	return dump
}

func (dump *Dump) Lookup(key string) interface{} {
	val, exist := dump.dmap[key]
	if !exist {
		// TODO 查找本地文件
	}
	return val
}

func (dump *Dump) Extend(key string) {
	dump.Lock()
	dump.dmap[key] = struct{}{}
	dump.Unlock()
}
