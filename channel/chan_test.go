package channel_test

import (
	"import.moetang.info/go/nekoq-common/channel"
	"testing"
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
