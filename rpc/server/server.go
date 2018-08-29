package server

import "context"

type Server interface {
	Handle(serviceName string, handler Handler)
}

type Handler func(req interface{}, ctx context.Context) (resp interface{}, err error)
