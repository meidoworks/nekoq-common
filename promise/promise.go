package promise

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	TIMEOUT = errors.New("promise: await timeout")
)

type Promise interface {
	AddListener(l Listener) Promise

	IsDone() bool
	IsFailed() bool

	Get() (interface{}, error)
	//TODO	GetOrDefault(interface{}) (interface{}, error)

	Await() error
	//TODO	AwaitWithin(time.Duration) error
}

type Listener interface {
	OnResult(i interface{})
	OnError(e error)
}

type SettablePromise interface {
	Promise

	DoneWith(i interface{}) SettablePromise
	FailWith(e error) SettablePromise
}

type promiseImpl struct {
	sync.Once
	sync.Mutex
	sync.WaitGroup

	l *list.List

	done   bool
	fail   bool
	finish int32

	o interface{}
	e error
}

func NewSettablePromise(useListener bool) SettablePromise {
	promise := &promiseImpl{}
	if useListener {
		promise.l = list.New()
	}
	promise.Add(1)
	promise.finish = 0
	return promise
}

func (p *promiseImpl) DoneWith(i interface{}) SettablePromise {
	p.Do(func() {
		p.o = i
		p.done = true
		atomic.AddInt32(&p.finish, 1)
		isLocked := false
		if p.l != nil {
			p.Lock()
			isLocked = true
			for e := p.l.Front(); e != nil; e = e.Next() {
				e.Value.(Listener).OnResult(i)
			}
		}
		p.Done()
		if isLocked {
			p.Unlock()
		}
	})
	return p
}

func (p *promiseImpl) FailWith(e error) SettablePromise {
	p.Do(func() {
		p.e = e
		p.fail = true
		atomic.AddInt32(&p.finish, 1)
		isLocked := false
		if p.l != nil {
			p.Lock()
			isLocked = true
			for elem := p.l.Front(); elem != nil; elem = elem.Next() {
				elem.Value.(Listener).OnError(e)
			}
		}
		p.Done()
		if isLocked {
			p.Unlock()
		}
	})
	return p
}

func (p *promiseImpl) IsDone() bool {
	return p.done
}

func (p *promiseImpl) IsFailed() bool {
	return p.fail
}

func (p *promiseImpl) Get() (interface{}, error) {
	p.Await()
	return p.o, p.e
}

func (p *promiseImpl) AddListener(l Listener) Promise {
	if p.l == nil {
		return p
	}
	p.Lock()
	if atomic.LoadInt32(&p.finish) == 1 {
		if p.IsDone() {
			l.OnResult(p.o)
		} else if p.IsFailed() {
			l.OnError(p.e)
		} else {
			p.Unlock()
			panic("should not reach here.")
		}
	} else {
		p.l.PushBack(l)
	}
	p.Unlock()
	return p
}

func (p *promiseImpl) Await() error {
	p.WaitGroup.Wait()
	return nil
}

type listenerImpl struct {
	d func(interface{})
	f func(error)
}

func NewListener(resultFunc func(interface{}), errorFunc func(error)) Listener {
	return &listenerImpl{
		d: resultFunc,
		f: errorFunc,
	}
}

func (l *listenerImpl) OnResult(r interface{}) {
	l.d(r)
}

func (l *listenerImpl) OnError(e error) {
	l.f(e)
}
