package file

import (
	"errors"
	"os"
)

const (
	MAX_FALLOC_FILE_SIZE = (int64(1) << 31) - 1
)

// less than 2G file, allocation
func Falloc(filePath string, fileSize int64) error {
	err := falloc(fileSize, filePath)
	return err
}

func checkParam(fileSize int64, filePath string) error {
	if fileSize < 1 {
		return errors.New("file size < 1")
	}
	if fileSize > MAX_FALLOC_FILE_SIZE {
		return errors.New("file size exceeds (2^31 - 1)")
	}
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return errors.New("file exists. cannot allocate.")
	}
	return nil
}
