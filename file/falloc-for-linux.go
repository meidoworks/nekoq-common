// +build linux

package file

/*
#include <fcntl.h>
*/
import "C"

import (
	"os"
	"syscall"
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

	errno := C.posix_fallocate(C.int(f.Fd()), 0, C.off_t(int32(fileSize)))
	if errno != 0 {
		return syscall.Errno(errno)
	}

	return nil
}
