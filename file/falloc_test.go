package file_test

import (
	"testing"
)

import (
	"import.moetang.info/go/nekoq-common/file"
)

func TestPreAllocateFile(t *testing.T) {
	file.Falloc("./file.alloc", 54321)
}
