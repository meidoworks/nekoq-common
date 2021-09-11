package file_test

import (
	"testing"

	"github.com/meidoworks/nekoq-common/file"
)

func TestPreAllocateFile(t *testing.T) {
	file.Falloc("./file.alloc", 54321)
}
