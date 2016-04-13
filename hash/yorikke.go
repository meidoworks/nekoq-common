package hash

import (
	"unsafe"
)

// fnv1a_hash_yorikke derivative of FNV1A from [http://www.sanmayce.com/Fastest_Hash/]
func Fnv1a_hash_yorikke(p uintptr, wrdlen int) uint32 {
	var hash32 uint32 = 2166136261
	var hash32b uint32 = 2166136261

	for wrdlen >= 2*2*(4) {
		hash32 = (hash32 ^ (_rotl_KAZE(*(*uint32)(unsafe.Pointer(p))) ^ *(*uint32)(unsafe.Pointer(p + 4)))) * _YORIKKE_PRIME
		hash32b = (hash32b ^ (_rotl_KAZE(*(*uint32)(unsafe.Pointer(p + 8))) ^ *(*uint32)(unsafe.Pointer(p + 12)))) * _YORIKKE_PRIME
		wrdlen -= 2 * 2 * (4)
		p += 2 * 2 * (4)
	}

	// cases: 0..15
	if (wrdlen & 2 * (4)) != 0 {
		hash32 = (hash32 ^ *(*uint32)(unsafe.Pointer(p))) * _YORIKKE_PRIME
		hash32b = (hash32b ^ *(*uint32)(unsafe.Pointer(p + 4))) * _YORIKKE_PRIME
		p += 4 * (2)
	}
	// cases: 0..7
	if wrdlen&(4) != 0 {
		hash32 = (hash32 ^ uint32(*(*uint16)(unsafe.Pointer(p)))) * _YORIKKE_PRIME
		hash32b = (hash32b ^ uint32(*(*uint16)(unsafe.Pointer(p + 2)))) * _YORIKKE_PRIME
		p += 4 * (2)
	}
	if wrdlen&(2) != 0 {
		hash32 = (hash32 ^ uint32(*(*uint16)(unsafe.Pointer(p)))) * _YORIKKE_PRIME
		p += (2)
	}
	if wrdlen&1 != 0 {
		hash32 = (hash32 ^ uint32(*(*uint8)(unsafe.Pointer(p)))) * _YORIKKE_PRIME
	}

	hash32 = (hash32 ^ _rotl_KAZE(hash32b)) * _YORIKKE_PRIME

	return hash32 ^ (hash32 >> 16)
}

func _rotl_KAZE(x uint32) uint32 {
	return (((x) << 5) | ((x) >> (32 - 5)))
}
