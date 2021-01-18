package cache

type DumpFile interface {
	Lookup(key string) interface{}
	Extend(key string)
}

type DumpMap map[string]interface{}

func MakeDumpMap(pipe chan string) DumpMap {
	var dmap = make(DumpMap)
	go func() {
		for {
			select {
			case key := <-pipe:
				dmap[key] = struct{}{}
			}
		}
	}()
	return dmap
}

func (dump DumpMap) Lookup(key string) interface{} {
	val, exist := dump[key]
	if !exist {
		return nil
	}
	return val
}

func (dump DumpMap) Extend(pipe chan string, key string) {
	go func() { pipe <- key }()
}
