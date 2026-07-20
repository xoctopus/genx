package contextx

// T simple
// +genx:ctx
type T interface {
}

// TT simple
// +genx:ctx
type TT int

// Targs generic will be skipped
// +genx:ctx
type Targs[T any] struct{}

// Alias will be skipped
// +genx:ctx
type Alias = TT
