package channel_test

import (
	"testing"
)

import (
	"github.com/meidoworks/nekoq-common/channel"
)

func TestJudgeSendToClosedChannel(t *testing.T) {
	defer channel.JudgeSendToClosedChannel(channel.OMIT_SEND_TO_CLOSED_CHANNEL_ERROR())
	err := funcToTest()
	t.Log(err)
}

func funcToTest() (resultError error) {
	ch := make(chan bool)
	close(ch)
	defer channel.JudgeSendToClosedChannel(func(err error) {
		resultError = err
	})
	ch <- true
	return
}

func BenchmarkJudgeSendToClosedChannelWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcWithDefer()
	}
}
func BenchmarkJudgeSendToClosedChannelWithDeferNop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcWithDeferNop()
	}
}
func BenchmarkJudgeSendToClosedChannelWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcWithoutDefer()
	}
}

func funcWithDefer() (resultError error) {
	ch := make(chan bool, 1)
	defer channel.JudgeSendToClosedChannel(func(err error) {
		resultError = err
	})
	ch <- true
	return
}

func funcWithDeferNop() (resultError error) {
	ch := make(chan bool, 1)
	defer channel.JudgeSendToClosedChannel(channel.OMIT_SEND_TO_CLOSED_CHANNEL_ERROR())
	ch <- true
	return
}

func funcWithoutDefer() (resultError error) {
	ch := make(chan bool, 1)
	ch <- true
	return
}
