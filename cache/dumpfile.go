package cache

type DumpFile interface {
	Lookup(key string) interface{}
	Extend(key string)
}

type Dump struct {
	pipe chan string
	dmap map[string]interface{}
}

func NewDump() *Dump {
	dump := new(Dump)
	dump.pipe = make(chan string)
	dump.dmap = make(map[string]interface{})
	go func() {
		for {
			select {
			case key := <-dump.pipe:
				dump.dmap[key] = struct{}{}
			}
		}
	}()
	return dump
}

func (dump *Dump) Lookup(key string) interface{} {
	val, exist := dump.dmap[key]
	if !exist {
		return nil
	}
	return val
}

func (dump *Dump) Extend(key string) {
	go func() { dump.pipe <- key }()
}
