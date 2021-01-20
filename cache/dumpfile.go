package cache

import "sync"

type DumpFile interface {
	Lookup(string) interface{}
	Extend(string, func() error) error
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
	val, _ := dump.dmap[key]
	return val
}

func (dump *Dump) Extend(key string, run func() error) error {
	dump.Lock()
	defer dump.Unlock()

	if dump.Lookup(key) != nil {
		return nil
	}

	err := run()
	if err == nil {
		dump.dmap[key] = struct{}{}
	}
	return err
}
