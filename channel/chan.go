package channel

import (
	"errors"
	"reflect"
)

func omitSendToClosedChannelError(err error) {
}

func OMIT_SEND_TO_CLOSED_CHANNEL_ERROR() func(error) {
	return omitSendToClosedChannelError
}

/**
func Usage() (resultErr error) {
	defer channel.JudgeSendToClosedChannel(func(err error) {
		resultError = err
	})
	// send to channel
	return
}
*/
func JudgeSendToClosedChannel(f func(error)) {
	recovered := recover()
	if recovered != nil {
		if reflect.TypeOf(recovered).String() == "string" && recovered.(string) == "send on closed channel" {
			resultErr := errors.New("send on closed channel")
			f(resultErr)
		} else {
			panic(recovered)
		}
	} else {
		return
	}
}
