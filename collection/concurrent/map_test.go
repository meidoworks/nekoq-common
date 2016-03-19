package concurrent

import (
	"testing"
)

type item struct {
	hashCode int32
	name     string
}

func (this *item) Hash() int32 {
	return this.hashCode
}

func TestUsage(t *testing.T) {
	i1 := &item{1, "item1"}
	i2 := &item{2, "item2"}
	i11 := &item{1, "item1"}

	m := NewMap(CL_256)
	t.Log(m.Get(i1))
	t.Log(m.Put(i1, "value1"))
	t.Log(m.Put(i1, "value11"))
	t.Log(m.Get(i2))
	t.Log(m.Put(i2, "value2"))
	t.Log(m.Get(i1), m.Get(i2), m.Get(i11))
	t.Log(m.Remove(i1))
	t.Log(m.Get(i1), m.Get(i2), m.Get(i11))
}
