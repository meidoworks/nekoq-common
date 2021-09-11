package idgen

import "testing"

func TestNewIdGen(t *testing.T) {
	gen := NewIdGen(1, 1)

	for i := 0; i < 1000; i++ {
		t.Log(gen.GenerateString())
	}
}

func TestNewIdGen2(t *testing.T) {
	gen := NewIdGen(1, 1)

	for i := 0; i < 1000; i++ {
		t.Log(gen.Generate())
	}
}

func BenchmarkNewIdGen(b *testing.B) {
	gen := NewIdGen(1, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.GenerateString()
	}
}
