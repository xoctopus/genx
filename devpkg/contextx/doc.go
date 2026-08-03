// Package contextx provides a generator for creating type-safe context injection and extraction functions.
//
// # How to Integrate Generator
//
//	import _ "github.com/xoctopus/genx/devpkg/contextx"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...)
//
// # Code Definition Conventions
//
// The generator supports both normal types and generic types. It will skip aliases.
//
//  1. Normal Types:
//     Add `// +genx:context` to a struct, interface, or basic type.
//     The generator will create `From`, `With`, `Must`, and `Carry` functions
//     specific to that type.
//
//  2. Generic Types:
//     Add `// +genx:context` to an uninstantiated generic struct or interface.
//     The generator will create generic versions of the context functions,
//     preserving the type parameters and their constraints.
//
// Example:
//
//	package example
//
//	import (
//		"iter"
//		"net"
//	)
//
//	// User represents a normal type.
//	// +genx:context
//	type User struct {
//		Name string
//	}
//
//	// GenericType represents a generic type with constraints.
//	// +genx:context
//	type GenericType[T any] struct {
//		V T
//	}
//
// Generated Code (Conceptual):
//
//	// For User:
//	func UserFrom(ctx context.Context) (User, bool) { ... }
//	func MustUser(ctx context.Context) User { ... }
//	func WithUser(ctx context.Context, v User) context.Context { ... }
//	func CarryUser(v User) contextx.Carrier { ... }
//
//	// For GenericType:
//	func GenericTypeFrom[T any](ctx context.Context) (GenericType[T], bool) { ... }
//	func MustGenericType[T any](ctx context.Context) GenericType[T] { ... }
//	func WithGenericType[T any](ctx context.Context, v GenericType[T]) context.Context { ... }
//	func CarryGenericType[T any](v GenericType[T]) contextx.Carrier { ... }
package contextx
