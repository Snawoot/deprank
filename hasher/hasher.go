package hasher

import (
	"github.com/dolthub/maphash"
)

type Hasher[K comparable] maphash.Hasher[K]

func NewHasher[K comparable]() Hasher[K] {
	return Hasher[K](maphash.NewHasher[K]())
}

func (h Hasher[K]) Hash(key K) uint32 {
	return uint32(maphash.Hasher[K](h).Hash(key))
}

func (h Hasher[K]) Equal(a, b K) bool {
	return a == b
}
