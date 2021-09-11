package idgen

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/meidoworks/nekoq-common/utilities"
)

var defaultGen = newIdGen(0, 0)

type nekoqIdGen struct {
	lock sync.Mutex

	time    int64
	curTime time.Time
	seq     int32

	nodeIdMask    int64
	elementIdMask int64

	workerStarter *sync.Once
}

func (n *nekoqIdGen) Generate() [2]int64 {
	id, err := n.GenerateWithError()
	if err != nil {
		panic(err)
	}
	return id
}

func (n *nekoqIdGen) GenerateString() string {
	id := n.Generate()
	return utilities.Int64Hex(id[0]) + utilities.Int64Hex(id[1])
}

func (n *nekoqIdGen) GenerateWithError() ([2]int64, error) {
	return n.Next()
}

func (n *nekoqIdGen) GenerateStringWithError() (string, error) {
	id, err := n.GenerateWithError()
	if err != nil {
		return "", err
	}
	return utilities.Int64Hex(id[0]) + utilities.Int64Hex(id[1]), nil
}

var _ IdGen128 = new(nekoqIdGen)

var (
	_EMPTY_RESULT         = [2]int64{0, 0}
	_EMPTY_RANGE_RESULT   = [][2]int64{}
	_ERROR_CLOCK_BACKWARD = errors.New("clock backward")
	_MAX_VALUE_INT32      = int32(0x7fffffff)
)

const (
	_START_TIME_MILLIS int64 = 1521639000000 // 20180321213000
)

// 48 bits time + 16 bits nodeId + 32 bits elementId + 32 bits inc
func newIdGen(nodeId int16, elementId int32) *nekoqIdGen {
	gen := &nekoqIdGen{
		time:          0,
		seq:           0,
		nodeIdMask:    int64(nodeId) & 0x000000000000FFFF,
		elementIdMask: (int64(elementId) & 0x00000000FFFFFFFF) << 32,
		workerStarter: new(sync.Once),
	}

	gen.workerStarter.Do(gen.startTimeWorker)

	return gen
}

func (this *nekoqIdGen) startTimeWorker() {
	now := time.Now()
	this.curTime = now
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			t := <-ticker.C
			this.curTime = t
		}
	}()
}

func (id *nekoqIdGen) getTimeMillis() int64 {
	n := id.curTime
	return int64((n.Unix()*1000 + (int64(n.Nanosecond()%1000000000) / 1000000)) & 0x7fffffffffffffff)
}

func (this *nekoqIdGen) NextN(cnt int) ([][2]int64, error) {
	result := make([][2]int64, cnt)
	retryCnt := 0

	this.lock.Lock()
	for {
		timeInMills := this.getTimeMillis()

		if timeInMills > this.time {
			// set seq to zero & return result
			this.time = timeInMills
			this.seq = int32(cnt - 1)
			this.lock.Unlock()
			return makeIdRange(timeInMills, this.nodeIdMask, this.elementIdMask, result, 0, int32(cnt-1)), nil
		} else if timeInMills == this.time {
			newSeq := this.seq + int32(cnt-1)
			// inc seq or wait until next time
			if newSeq < _MAX_VALUE_INT32 {
				// inc seq
				prevSeq := this.seq
				this.seq = newSeq
				this.lock.Unlock()
				return makeIdRange(timeInMills, this.nodeIdMask, this.elementIdMask, result, prevSeq, newSeq-1), nil
			} else {
				// wait until next time
				newtime := this.tillNextMillisecond(timeInMills)
				// success
				this.time = newtime
				this.seq = int32(cnt - 1)
				this.lock.Unlock()
				return makeIdRange(newtime, this.nodeIdMask, this.elementIdMask, result, 0, int32(cnt-1)), nil
			}
		} else {
			if retryCnt < 3 {
				retryCnt++
				continue
			}

			// error: clock backward
			this.lock.Unlock()
			return _EMPTY_RANGE_RESULT, _ERROR_CLOCK_BACKWARD
		}
	}
}

func (id *nekoqIdGen) Next() ([2]int64, error) {
	retryCnt := 0

	id.lock.Lock()
	for {
		timeInMills := id.getTimeMillis()

		if timeInMills > id.time {
			// set seq to zero & return result
			id.time = timeInMills
			id.seq = 0
			id.lock.Unlock()
			return makeId(timeInMills, id.nodeIdMask, id.elementIdMask, 0), nil
		} else if timeInMills == id.time {
			// inc seq or wait until next time
			if id.seq < _MAX_VALUE_INT32 {
				// inc seq
				id.seq = id.seq + 1
				newseq := id.seq
				id.lock.Unlock()
				return makeId(timeInMills, id.nodeIdMask, id.elementIdMask, newseq), nil
			} else {
				// wait until next time
				newtime := id.tillNextMillisecond(timeInMills)
				// success
				id.time = newtime
				id.seq = 0
				id.lock.Unlock()
				return makeId(newtime, id.nodeIdMask, id.elementIdMask, 0), nil
			}
		} else {
			if retryCnt < 3 {
				retryCnt++
				continue
			}

			// error: clock backward
			id.lock.Unlock()
			return _EMPTY_RESULT, _ERROR_CLOCK_BACKWARD
		}
	}
}

func makeIdRange(time, nodeIdMask int64, elementId int64, result [][2]int64, seqStart int32, seqEnd int32) [][2]int64 {
	for idx, start := 0, seqStart; start <= seqEnd; idx, start = idx+1, start+1 {
		l := elementId | (int64(start) & 0x00000000ffffffff)
		result[idx] = [2]int64{((time - _START_TIME_MILLIS) << 16) | nodeIdMask, l}
	}
	return result
}

func makeId(time, nodeIdMask int64, elementId int64, seq int32) [2]int64 {
	l := elementId | (int64(seq) & 0x00000000ffffffff)
	return [2]int64{((time - _START_TIME_MILLIS) << 16) | nodeIdMask, l}
}

func (id *nekoqIdGen) tillNextMillisecond(time int64) int64 {
	for {
		newtime := id.getTimeMillis()
		if newtime > time {
			return newtime
		}
		runtime.Gosched()
	}
}
