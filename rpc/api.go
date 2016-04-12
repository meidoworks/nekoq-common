package rpc

import (
	"moetang.info/go/common/async"
	"moetang.info/go/common/context"
)

type Client interface {
	CallSync(method string, param interface{}, appInfo *context.AppInfo) (interface{}, error)
	CallAsync(method string, param interface{}, AppInfo *context.AppInfo) (async.Future, error)

	CloseSync() error
}
