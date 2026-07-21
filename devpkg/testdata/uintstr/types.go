package uintstr

// OrderID
// +genx:uintstr
type OrderID uint64

// AliasOrderID alias will be skipped
// +genx:uintstr
type AliasOrderID = OrderID

// Generic will be skipped
// +genx:uintstr
type Generic[T any] struct{}

// ItemID not unsigned will be skipped
// +genx:uintstr
type ItemID int64

// CodeID not basic will be skipped
// +genx:uintstr
type CodeID []byte
