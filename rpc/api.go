package rpc

import (
	"import.moetang.info/go/nekoq-common/async"
	"import.moetang.info/go/nekoq-common/context"
)

type Client interface {
	CallSync(method string, param interface{}, appInfo *context.AppInfo) (interface{}, error)
	CallAsync(method string, param interface{}, AppInfo *context.AppInfo) (async.Future, error)

	CloseSync() error
}
