package utilities

import "testing"

func TestInt64Hex(t *testing.T) {
	t.Log(Int64Hex(9199999999999999999))
	t.Log(Int64Hex(-1))
	t.Log(Int64Hex(4357640193405743614))
	t.Log(Int64Hex(-2))
	t.Log(Int64Hex(1))
}
