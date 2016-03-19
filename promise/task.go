package promise

import (
	"errors"
	"sync/atomic"
	"time"
)

//================================================================
// basic infrastructure
//================================================================

//---------------- consts ----------------

const (
	PRIORITY_MAX    int32 = (int32)(((int64(1)) << 31) - 1)
	PRIORITY_MIN    int32 = 1
	PRIORITY_NORMAL int32 = 65535

	MAX_TASKS_PER_PARSEQ = 100000
)

var (
	g_PREDICATE_FAILED_ERROR = errors.New("promise: predicate failed.")
	g_EMPTY_RESULT_ERROR     = errors.New("promise: no result.")
)

var (
	g_ILLEGAL_ARGUMENT_PANIC = errors.New("promise: illegal argument panic.")
)

var (
	g_EMPTY_INPUT = &input{}
)

func ERROR_PREDICATE_FAILED() error {
	return g_PREDICATE_FAILED_ERROR
}
func ERROR_EMTPY_RESULT() error {
	return g_EMPTY_RESULT_ERROR
}
func PANIC_ILLEGAL_ARUGUMENT() error {
	return g_ILLEGAL_ARGUMENT_PANIC
}

//---------------- apis ----------------

type Input interface {
	GetResult() Result
}

type ParOutput interface {
	GetResults() []Result
	//TODO	GetResult(string) Task
}
type ParError interface {
	ParOutput

	Error() string
}

type TaskInfo interface {
	GetName() string
	GetPriority() int32
	IsRequired() bool
}

type Result interface {
	TaskInfo

	Result() Promise
}

type Task interface {
	TaskInfo

	Do() (Promise, Context)

	SetPriority(int32)

	//if non required task failed, task will go on processing.
	SetRequired(bool)

	SetPredicate(func(Input) bool)
	ClearPredicate()
}

type Trace interface {
	//TODO
}

type Context interface {
	GetTrace() Trace
}

//================================================================
// implementions of infrastructure
//================================================================

type context struct {
	SettablePromise
}

func (this *context) GetTrace() Trace {
	//TODO
	return nil
}

type task struct {
	t         func(Input) Promise
	name      string
	priority  int32
	required  bool
	predicate func(Input) bool
}

type input struct {
	input Result
}

type parInput struct {
	requiredCnt      int32
	requiredFinished int32
	requiredError    int32

	inputs   []Result
	errorMsg error
}

type inputItem struct {
	TaskInfo
	result Promise
}

func (this *parInput) Error() string {
	return this.errorMsg.Error()
}

func (this *parInput) GetResults() []Result {
	return this.inputs
}

func (this *inputItem) Result() Promise {
	return this.result
}

func (i *input) GetResult() Result {
	return i.input
}

func (this *task) GetName() string {
	return this.name
}
func (this *task) GetPriority() int32 {
	return this.priority
}
func (this *task) SetPriority(priority int32) {
	this.priority = priority
}
func (this *task) IsRequired() bool {
	return this.required
}
func (this *task) SetRequired(required bool) {
	this.required = required
}
func (this *task) SetPredicate(predicate func(Input) bool) {
	this.predicate = predicate
}
func (this *task) ClearPredicate() {
	this.predicate = nil
}

func (this *task) Do() (Promise, Context) {
	return runFunc(this, this.t, g_EMPTY_INPUT)
}

func (this *task) do(i Input) (Promise, Context) {
	return runFunc(this, this.t, i)
}

//================================================================
// baisc api
//================================================================

// promise must be settable
func NewTask(name string, priority int32, required bool, predicate func(Input) bool, f func(Input) Promise) Task {
	return newTask("u", name, priority, required, predicate, f)
}

// taskPrefix - u: user task, s: system task like par/seq
func newTask(taskPrefix string, name string, priority int32, required bool, predicate func(Input) bool, f func(Input) Promise) Task {
	ta := &task{
		name:      taskPrefix + "::" + name,
		priority:  priority,
		predicate: predicate,
		required:  required,
		t:         f,
	}
	return ta
}

func NewGoroutineTask(name string, priority int32, required bool, predicate func(Input) bool, f func(Input) (interface{}, error)) Task {
	t := func(input Input) Promise {
		promise := NewSettablePromise(true)
		go func() {
			r, e := f(input)
			if e != nil {
				promise.FailWith(e)
			} else {
				promise.DoneWith(r)
			}
		}()
		return promise
	}
	return NewTask(name, priority, required, predicate, t)
}

// work with predicate, reture last accepted result with come across first false
// ensure invoke order
// done type - interface{} / fail type - error
func Seq(tasks ...Task) Task {
	if tasks == nil || len(tasks) < 1 || len(tasks) > MAX_TASKS_PER_PARSEQ {
		panic(PANIC_ILLEGAL_ARUGUMENT())
	}
	f := func(input Input) Promise {
		p := NewSettablePromise(true)
		seqNext(p, 0, tasks, input)
		return p
	}
	return newTask("s", "seq_task", PRIORITY_NORMAL, true, nil, f)
}

// ensure all tasks with and only with required mark are done/error
// done type - ParOutput / fail type - ParError
func Par(tasks ...Task) Task {
	if tasks == nil || len(tasks) < 1 || len(tasks) > MAX_TASKS_PER_PARSEQ {
		panic(PANIC_ILLEGAL_ARUGUMENT())
	}
	f := func(input Input) Promise {
		p := NewSettablePromise(true)
		parProcess(p, tasks, input)
		return p
	}
	return newTask("s", "par_task", PRIORITY_NORMAL, true, nil, f)
}

//================================================================
// internal of par/seq
//================================================================

func parProcess(p SettablePromise, tasks []Task, i Input) {
	result := &parInput{}
	result.inputs = make([]Result, len(tasks))
	result.requiredFinished = 0
	result.requiredError = 0

	for _, t := range tasks {
		required := t.IsRequired()
		if required {
			result.requiredCnt = result.requiredCnt + 1
		}
	}

	for idx, t := range tasks {
		rp, _ := t.(*task).do(i) //must use buildin type

		ii := &inputItem{}
		ii.result = rp
		ii.TaskInfo = t
		result.inputs[idx] = ii

		required := t.IsRequired()
		l := NewListener(func(i interface{}) {
			if !required {
				return
			}
			newFinished := atomic.AddInt32(&result.requiredFinished, 1)
			if newFinished == result.requiredCnt {
				p.DoneWith(ParOutput(result))
			}
		}, func(e error) {
			if !required {
				return
			}
			result.errorMsg = e
			p.FailWith(ParError(result))
		})

		rp.AddListener(l)
	}
}

func seqProcess(p SettablePromise, t Task, idx int, tasks []Task, i Input) {
	rp, _ := t.(*task).do(i) //must use buildin type
	ii := &inputItem{}
	ii.result = rp
	ii.TaskInfo = t
	toInput := &input{
		input: ii,
	}
	l := NewListener(func(i interface{}) {
		seqNext(p, idx+1, tasks, toInput)
	}, func(e error) {
		if e == ERROR_PREDICATE_FAILED() {
			if i.GetResult() == nil {
				p.FailWith(g_EMPTY_RESULT_ERROR)
				return
			}
			finalPromise := i.GetResult().Result()
			doneItem, errorItem := finalPromise.Get()
			if finalPromise.IsDone() {
				p.DoneWith(doneItem)
			} else {
				p.FailWith(errorItem)
			}
		} else {
			seqNext(p, idx+1, tasks, toInput)
		}
	})
	rp.AddListener(l)
}
func seqNext(p SettablePromise, nextIdx int, tasks []Task, i Input) {
	if nextIdx < len(tasks) {
		seqProcess(p, tasks[nextIdx], nextIdx, tasks, i)
	} else if nextIdx == len(tasks) {
		finalPromise := i.GetResult().Result()
		doneItem, errorItem := finalPromise.Get()
		if finalPromise.IsDone() {
			p.DoneWith(doneItem)
		} else {
			p.FailWith(errorItem)
		}
	}
}

//================================================================
// internal utils
//================================================================

func runFunc(t *task, f func(Input) Promise, i Input) (Promise, Context) {
	td := &timeDuraion{}
	td.task = t
	sp := NewSettablePromise(true)
	c := &context{}
	c.SettablePromise = sp

	td.start = time.Now().UnixNano()
	var p Promise
	if t.predicate != nil && !t.predicate(i) {
		p = NewSettablePromise(true).FailWith(ERROR_PREDICATE_FAILED())
	} else {
		p = f(i)
	}
	l := NewListener(func(r interface{}) {
		td.end = time.Now().UnixNano()
		c.DoneWith(td)
	}, func(e error) {
		td.end = time.Now().UnixNano()
		c.DoneWith(td)
	})
	p.AddListener(l)

	return p, c
}

type timeDuraion struct {
	task  Task
	start int64
	end   int64
}
