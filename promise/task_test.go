package promise

import (
	"errors"
	"testing"
)

func ExampleTaskUsage() {
	Par()
	var _ Task = &task{}
}

func TestUsage(t *testing.T) {
	Seq(&task{})
	Par(&task{})
}

func TestUsageTask(t *testing.T) {
	task := NewTask("hello",
		1,
		true,
		func(input Input) bool {
			t.Log("run predicate")
			return true
		},
		func(input Input) Promise {
			t.Log("do work.")
			return NewSettablePromise(true).DoneWith(100)
		})
	p, c := task.Do()
	t.Log(p.Get())
	t.Log(c.GetTrace())

	_, err := p.Get()
	t.Log(err == ERROR_PREDICATE_FAILED(), err == errors.New(ERROR_PREDICATE_FAILED().Error()))
}

func TestUsageGoroutineTask(t *testing.T) {
	task := NewGoroutineTask("hello2",
		1,
		true,
		nil,
		func(input Input) (interface{}, error) {
			t.Log("do work2.")
			return "return_result", nil
		})
	p, c := task.Do()
	t.Log(p.Get())
	t.Log(c.GetTrace())
}

func BenchmarkTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := NewTask("hello",
			1,
			true,
			func(input Input) bool {
				return true
			},
			func(input Input) Promise {
				return NewSettablePromise(true).DoneWith("return_result")
			})
		task.Do()
	}
}

func BenchmarkTaskSingleton(b *testing.B) {
	task := NewTask("hello",
		1,
		true,
		func(input Input) bool {
			return true
		},
		func(input Input) Promise {
			return NewSettablePromise(true).DoneWith("return_result")
		})
	for i := 0; i < b.N; i++ {
		task.Do()
	}
}

func BenchmarkGoroutineTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := NewGoroutineTask("hello",
			1,
			true,
			func(input Input) bool {
				return true
			},
			func(input Input) (interface{}, error) {
				return "return_result", nil
			})
		task.Do()
	}
}

func BenchmarkGoroutineTaskSingleton(b *testing.B) {
	task := NewGoroutineTask("hello",
		1,
		true,
		func(input Input) bool {
			return true
		},
		func(input Input) (interface{}, error) {
			return "return_result", nil
		})
	for i := 0; i < b.N; i++ {
		task.Do()
	}
}

func TestSeqTasks(t *testing.T) {
	tttt := NewTask("hello",
		1,
		true,
		func(i Input) bool {
			return true
		},
		func(i Input) Promise {
			t.Log("input:", i)
			var appendStr string
			if i.(*input).input == nil {
			} else {
				pr, _ := i.(*input).input.Result().Get()
				appendStr = pr.(string)
			}
			return NewSettablePromise(true).DoneWith("return_result" + " " + appendStr)
		})

	seqTask := Seq(tttt, tttt, tttt, tttt)

	p, c := seqTask.Do()
	t.Log(p.Get())
	t.Log(c.GetTrace())
}

func BenchmarkSeqTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tttt := NewTask("hello",
			1,
			true,
			func(i Input) bool {
				return true
			},
			func(i Input) Promise {
				return NewSettablePromise(true).DoneWith("return_result")
			})

		seqTask := Seq(tttt)

		seqTask.Do()
	}
}

func BenchmarkSeqTaskSingleton(b *testing.B) {
	tttt := NewTask("hello",
		1,
		true,
		func(i Input) bool {
			return true
		},
		func(i Input) Promise {
			return NewSettablePromise(true).DoneWith("return_result")
		})

	seqTask := Seq(tttt)

	for i := 0; i < b.N; i++ {
		seqTask.Do()
	}
}

func TestParTask(t *testing.T) {

	tttt := NewTask("hello",
		1,
		true,
		func(i Input) bool {
			return true
		},
		func(i Input) Promise {
			return NewSettablePromise(true).DoneWith("return_result")
		})

	parTask := Par(tttt, tttt)

	p, _ := parTask.Do()

	r, e := p.Get()
	t.Log(len(r.(ParOutput).GetResults()), e)

	ttte := NewTask("helloerror",
		1,
		true,
		func(i Input) bool {
			return true
		},
		func(i Input) Promise {
			return NewSettablePromise(true).FailWith(errors.New("error~~~"))
		})

	parTask2 := Par(ttte, tttt)

	p, _ = parTask2.Do()
	r, e = p.Get()
	t.Log(r, e)
	t.Log(e.(ParError).GetResults()[0].Result().Get())
	t.Log(e.(ParError).GetResults()[1].Result().Get())
}

func BenchmarkParTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tttt := NewTask("hello",
			1,
			true,
			func(i Input) bool {
				return true
			},
			func(i Input) Promise {
				return NewSettablePromise(true).DoneWith("return_result")
			})

		parTask := Par(tttt)

		parTask.Do()
	}
}

func BenchmarkParTaskSingleton(b *testing.B) {
	tttt := NewTask("hello",
		1,
		true,
		func(i Input) bool {
			return true
		},
		func(i Input) Promise {
			return NewSettablePromise(true).DoneWith("return_result")
		})

	parTask := Par(tttt)

	for i := 0; i < b.N; i++ {
		parTask.Do()
	}
}

func BenchmarkParTaskForNew100Tasks(b *testing.B) {
	tasks := make([]Task, 100)

	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			tasks[i] = NewTask("hello",
				1,
				true,
				func(i Input) bool {
					return true
				},
				func(i Input) Promise {
					return NewSettablePromise(true).DoneWith("return_result")
				})
		}

		parTask := Par(Seq(tasks[72:86]...),
			Par(Seq(tasks[0:12]...),
				Seq(tasks[12:35]...),
				Seq(tasks[35:47]...),
				Seq(tasks[47:54]...),
				Seq(tasks[54:58]...),
				Seq(tasks[58:63]...),
				Seq(tasks[63:72]...)),
			Par(tasks[86:92]...),
			Seq(tasks[92:]...),
		)

		parTask.Do()
	}
}

func BenchmarkParTaskFor100Tasks(b *testing.B) {
	tasks := make([]Task, 100)

	for i := 0; i < 100; i++ {
		tasks[i] = NewTask("hello",
			1,
			true,
			func(i Input) bool {
				return true
			},
			func(i Input) Promise {
				return NewSettablePromise(true).DoneWith("return_result")
			})
	}

	parTask := Par(Seq(tasks[72:86]...),
		Par(Seq(tasks[0:12]...),
			Seq(tasks[12:35]...),
			Seq(tasks[35:47]...),
			Seq(tasks[47:54]...),
			Seq(tasks[54:58]...),
			Seq(tasks[58:63]...),
			Seq(tasks[63:72]...)),
		Par(tasks[86:92]...),
		Seq(tasks[92:]...),
	)

	for i := 0; i < b.N; i++ {
		parTask.Do()
	}
}

func BenchmarkParTaskForNew100TasksParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		tasks := make([]Task, 100)

		for pb.Next() {
			for i := 0; i < 100; i++ {
				tasks[i] = NewTask("hello",
					1,
					true,
					func(i Input) bool {
						return true
					},
					func(i Input) Promise {
						return NewSettablePromise(true).DoneWith("return_result")
					})
			}

			parTask := Par(Seq(tasks[72:86]...),
				Par(Seq(tasks[0:12]...),
					Seq(tasks[12:35]...),
					Seq(tasks[35:47]...),
					Seq(tasks[47:54]...),
					Seq(tasks[54:58]...),
					Seq(tasks[58:63]...),
					Seq(tasks[63:72]...)),
				Par(tasks[86:92]...),
				Seq(tasks[92:]...),
			)

			parTask.Do()
		}
	})
}
