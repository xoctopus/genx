package contextx

import (
	"net"
)

// T simple
// +genx:context
type T interface {
}

// TT simple
// +genx:context
type TT int

// Targs generic will be skipped
// +genx:context
type Targs[T any] struct{}

// Alias will be skipped
// +genx:context
type Alias = TT

// Generic
// +genx:context
type Generic[T net.Addr] struct{}

// // GenericRecursive TODO recursive type parameters is not unsupported
// // +genx:context
// type GenericRecursive[T any, S iter.Seq2[int, T]] struct{}
