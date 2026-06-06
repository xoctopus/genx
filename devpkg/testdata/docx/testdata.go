// Package docx define test cases for docx generation
// +genx:doc
package docx

import (
	"github.com/xoctopus/genx/devpkg/testdata/docx/ext"
)

var (
	// x integer
	x int
	// Y string
	Y string
)

// Stringer interface
type Stringer interface {
	String() string
}

// Alias alias
type Alias = Stringer

// Int union
type Int interface {
	int | int8
}

type Empty struct{}

// S struct
type S struct {
	// unexported
	unexported any
	// Empty EmptyPrefix
	Empty
	// Empty2 Empty2Doc
	Empty2 Empty
	// Inline InlineUnnamed
	Inline struct {
		A any // A
	}
	// NameType
	NameType int
	// Interface NamedInterface
	Interface Stringer
	// D AnonymousDPrefix
	D
	// S *SPrefix
	*S
	// T1 ext.T1Prefix
	ext.T1
	// T2 *ext.T2Prefix
	*ext.T2
	// T3 ext.T3[int]Prefix
	ext.T3[int]
	// T4 *ext.T4[int]Prefix
	*ext.T4[int]
	// T5 ext.T5[int,int]Prefix
	ext.T5[int, int]
	// T6 *ext.T6[int, int]Prefix
	*ext.T6[int, int]
	NoDoc any
}

// D doc
type D struct {
	// X int
	X int
	// Y string
	Y string
	// F func
	F F
}

// F func
type F func()

type G[V any] struct {
	// NameType
	NameType int
	// Interface
	Interface Stringer
	// D ident only
	D
	// S star expr
	*S
	// T1 ext selector
	ext.T1
	// T2 star ref
	*ext.T2
	// T3 index expr
	ext.T3[V]
	// T4 star index expr
	*ext.T4[V]
	// T5 index list expr
	ext.T5[V, V]
	// T6 star index list expr
	*ext.T6[V, V]
}

func A() {}

// FAlias alias of F
type FAlias = F
