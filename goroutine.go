package common

import (
	"math/rand"
)

type GoroutineContext interface {
	GetGoroutineRandomId() int64
}

type context struct {
	id int64
}

func (this *context) GetGoroutineRandomId() int64 {
	return this.id
}

func (this *context) Hash() int32 {
	return int32(this.id)
}

func NewGoroutineContext() GoroutineContext {
	c := &context{}
	c.id = rand.Int63()
	return c
}
