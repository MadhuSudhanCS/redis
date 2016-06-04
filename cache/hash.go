package cache

import (
	"hash"
	"hash/fnv"
)

var h hash.Hash64

func init() {
	h = fnv.New64a()
}

func Hash(data string) uint64 {
	h.Write([]byte(data))
	hv := h.Sum64()
	h.Reset()

	return hv
}
