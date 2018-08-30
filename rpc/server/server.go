package server

type Server interface {
	RegisterHandler(handler interface{}) error
}
