// +build !linux

package file

import (
	"os"
)

func falloc(fileSize int64, filePath string) error {
	err := checkParam(fileSize, filePath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	_, err = f.Seek(fileSize-1, 1)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte{0})
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
