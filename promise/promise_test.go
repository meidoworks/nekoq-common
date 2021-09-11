package promise

import (
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkCreationWithListenerSupport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSettablePromise(true)
	}
}
func BenchmarkCreationWithoutListenerSupport(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSettablePromise(false)
	}
}

func BenchmarkAtomicAdd(b *testing.B) {
	var ddd int64 = 0

	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&ddd, 1)
	}
}

func TestPromiseBasicUsage(t *testing.T) {
	p := NewSettablePromise(true)
	p.AddListener(NewListener(func(i interface{}) {
		t.Log(time.Now(), "pre listener finish")
	}, func(e error) {}))
	go func() {
		t.Log(time.Now(), "done")
		time.Sleep(5 * time.Second)
		p.DoneWith(1)
	}()
	t.Log(p.Get())
	t.Log(time.Now(), "finish")
	p.AddListener(NewListener(func(i interface{}) {
		t.Log(time.Now(), "listener finish")
	}, func(e error) {}))
}
