// Package testdata package level document
//
// comments for testdata package
package testdata

import "io"

// Structure is a struct type for testing
// line1
// line2
// +ignore=name
type Structure struct {
	// T same package
	// TODO
	// FIXME
	// UNDONE
	T any
}

/*
Fields for field document
this is a quoted multiline comments
*/
type Fields struct {
	// Name title
	// descriptions...
	V string
	// Structure has sub
	Structure
	// Named
	Named Structure
}

// GenDecl for type list
type (
	// A some type A
	A any
	// B some type B
	B int
	// C interface type
	C interface {
		io.Reader
	}
)

// GenDecl for var list
var (
	// VarA int
	VarA int
	// VarB float
	VarB float32
)

// GenDecl for const list
var (
	// ConstA int
	ConstA int
	// ConstB float
	ConstB float32
)
