package concurrent

import (
	"sync"

	"github.com/meidoworks/nekoq-api/object"
)

type ConcurrentMap interface {
	Put(key object.Object, value interface{}) interface{}
	Get(key object.Object) interface{}
	Remove(key object.Object) interface{}
}

type concurrentMap struct {
	concurrentLevel int32
	locks           []*sync.RWMutex
	maps            []map[object.Object]interface{}
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
	m.maps = make([]map[object.Object]interface{}, int(cl))

	for i := 0; i < int(cl); i++ {
		m.locks[i] = &sync.RWMutex{}
		m.maps[i] = make(map[object.Object]interface{})
	}

	return m
}

func (this *concurrentMap) Put(key object.Object, value interface{}) (preObj interface{}) {
	idx := int(key.HashCode() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.Lock()
	preObj = m[key]
	m[key] = value
	lock.Unlock()
	return
}

func (this *concurrentMap) Get(key object.Object) (obj interface{}) {
	idx := int(key.HashCode() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.RLock()
	obj = m[key]
	lock.RUnlock()
	return
}

func (this *concurrentMap) Remove(key object.Object) (preObj interface{}) {
	idx := int(key.HashCode() & this.concurrentLevel)
	lock := this.locks[idx]
	m := this.maps[idx]
	lock.Lock()
	preObj = m[key]
	delete(m, key)
	lock.Unlock()
	return
}
