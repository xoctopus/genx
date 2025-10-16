package testdata

func Func() {
Func:
	for i := range 100 {
		if i == 1 {
			goto Func
		}
	}
}

var Var int

var VarFuncT = FuncT[int]

type DemoEnum int8

const (
	DEMO_ENUM_UNKNOWN DemoEnum = iota
	DEMO_ENUM_A
	DEMO_ENUM_B
	DEMO_ENUM_C
)

type T[A any] struct{ v A }

func Func2() {
Loop:
	for i := range 100 {
		if i == 1 {
			goto Loop
		}
	}
}

func FuncT[T any]() {}

type Struct struct {
	A any
	B any
}

func (s *Struct) Func() {}

type StructT[T any, E any] struct {
	_ *StructT[int, string]
}

type String string
