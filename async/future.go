package async

type Future interface {
	Get() (interface{}, error)
}
