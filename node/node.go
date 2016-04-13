package node

import (
	"unsafe"
)

import (
	"moetang.info/go/common/hash"
)

type ItemId [2]int64

type NodeId ItemId

func (this ItemId) Hash() uint32 {
	return nodeIdHash(this)
}

func (this NodeId) Hash() uint32 {
	return nodeIdHash(this)
}

func nodeIdHash(data [2]int64) uint32 {
	p := uintptr(unsafe.Pointer(&data))

	return hash.Fnv1a_hash_yorikke(p, 16)
}
