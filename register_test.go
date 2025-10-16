package genx_test

import (
	"go/types"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx"
)

type g struct {
	name string
}

func (x *g) Identifier() string { return x.name }

func (x *g) Generate(_ genx.Context, _ types.Type) error { return nil }

func TestRegister(t *testing.T) {
	genx.Register(&g{name: "a"})
	genx.Register(&g{name: "b"})

	gs := genx.Get("a", "b", "b")
	Expect(t, len(gs), Equal(2))
	Expect(t, gs[0].Identifier(), Equal("a"))
	Expect(t, gs[1].Identifier(), Equal("b"))

	gs = genx.Get()
	Expect(t, len(gs), Equal(2))
	Expect(t, gs[0].Identifier(), Equal("a"))
	Expect(t, gs[1].Identifier(), Equal("b"))
}
