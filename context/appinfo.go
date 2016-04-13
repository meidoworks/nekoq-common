package context

type AppInfo struct {
}

func (this *AppInfo) Convert() ([]byte, error) {
	return []byte{}, nil
}

func FromBytes(data []byte) (*AppInfo, error) {
	return new(AppInfo), nil
}
