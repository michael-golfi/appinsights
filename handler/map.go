package handler

import "sync"

type logPairMap struct {
	sync.RWMutex
	internal map[string]*logPair
}

func newLogPairMap() *logPairMap {
	return &logPairMap{
		internal: make(map[string]*logPair),
	}
}

func (rm *logPairMap) Load(key string) (value *logPair, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *logPairMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *logPairMap) Store(key string, value *logPair) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}