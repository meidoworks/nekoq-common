package file_test

import (
	"testing"
)

import (
	"goimport.moetang.info/nekoq-common/file"
)

func TestPreAllocateFile(t *testing.T) {
	file.Falloc("./file.alloc", 54321)
}
