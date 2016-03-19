package concurrent

import (
	"sync"

	"moetang.info/go/common"
)

type ConcurrentMap interface {
	Put(key common.HashObj, value interface{}) interface{}
	Get(key common.HashObj) interface{}
	Remove(key common.HashObj) interface{}
}

type concurrentMap struct {
	concurrentLevel int32
	locks           []*sync.RWMutex
	maps            []map[common.HashObj]interface{}
}

type CONCURRENT_LEVEL int32

const (
	CL_1 CONCURRENT_LEVEL = (1)<<iota - 1
	CL_2
	CL_4
	CL_8
	CL_16
	CL_32
	CL_64
	CL_128
	CL_256
)

func NewMap(cl CONCURRENT_LEVEL) ConcurrentMap {
	m := &concurrentMap{}
	m.concurrentLevel = int32(cl)

	m.locks = make([]*sync.RWMutex, int(cl))
	m.maps = make([]map[common.HashObj]interface{}, int(cl))

	for i := 0; i < int(cl); i++ {
		m.locks[i] = &sync.RWMutex{}
		m.maps[i] = make(map[common.HashObj]interface{})
	}

	return m
}

func (this *concurrentMap) Put(key common.HashObj, value interface{}) (preObj interface{}) {
	idx := int(key.Hash() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.Lock()
	preObj = m[key]
	m[key] = value
	lock.Unlock()
	return
}

func (this *concurrentMap) Get(key common.HashObj) (obj interface{}) {
	idx := int(key.Hash() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.RLock()
	obj = m[key]
	lock.RUnlock()
	return
}

func (this *concurrentMap) Remove(key common.HashObj) (preObj interface{}) {
	idx := int(key.Hash() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.Lock()
	preObj = m[key]
	delete(m, key)
	lock.Unlock()
	return
}
