package client

import "context"

type Client interface {
	Send(req interface{}, ctx context.Context) (resp interface{}, err error)
}
