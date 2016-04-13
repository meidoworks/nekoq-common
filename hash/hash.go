package hash

/*
 *
 * this file is jesteress algorithm originally from apcera/gnatsd
 *
 */

import (
	"reflect"
	"unsafe"
)

type HashGenerator interface {
	Hash([]byte) uint32
}

// Jesteress derivative of FNV1A from [http://www.sanmayce.com/Fastest_Hash/]
func Jesteress(data []byte) uint32 {
	h32 := uint32(_OFF32)
	i, dlen := 0, len(data)

	for ; dlen >= _DDWSZ; dlen -= _DDWSZ {
		k1 := *(*uint64)(unsafe.Pointer(&data[i]))
		k2 := *(*uint64)(unsafe.Pointer(&data[i+8]))
		h32 = uint32((uint64(h32) ^ ((k1<<5 | k1>>27) ^ k2)) * _YP32)
		i += _DDWSZ
	}

	// Cases: 0,1,2,3,4,5,6,7
	if (dlen & _DWSZ) > 0 {
		k1 := *(*uint64)(unsafe.Pointer(&data[i]))
		h32 = uint32(uint64(h32)^k1) * _YP32
		i += _DWSZ
	}
	if (dlen & _WSZ) > 0 {
		k1 := *(*uint32)(unsafe.Pointer(&data[i]))
		h32 = (h32 ^ k1) * _YP32
		i += _WSZ
	}
	if (dlen & 1) > 0 {
		h32 = (h32 ^ uint32(data[i])) * _YP32
	}
	return h32 ^ (h32 >> 16)
}

func Yorikko(data []byte) uint32 {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	return Fnv1a_hash_yorikke(h.Data, h.Len)
}
