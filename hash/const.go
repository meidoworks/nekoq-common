package hash

// Constants for FNV1A and derivatives
const (
	_OFF32 = 2166136261
	_P32   = 16777619
	_YP32  = 709607
)

// Constants for multiples of sizeof(WORD)
const (
	_WSZ    = 4         // 4
	_DWSZ   = _WSZ << 1 // 8
	_DDWSZ  = _WSZ << 2 // 16
	_DDDWSZ = _WSZ << 3 // 32
)

// constants for fnv1a_hash_yorikke
const (
	_YORIKKE_PRIME uint32 = 709607
)
