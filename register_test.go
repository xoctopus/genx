package genx_test

import (
	"go/types"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/genx"
)

type g struct {
	name string
}

func (x *g) Identifier() string { return x.name }

func (x *g) Generate(_ genx.Context, _ types.Type) error { return nil }

func TestRegisterGenerator(t *testing.T) {
	genx.RegisterGenerator(&g{name: "a"})
	genx.RegisterGenerator(&g{name: "b"})

	gs := genx.GetGenerators("a", "b", "b")
	NewWithT(t).Expect(len(gs)).To(Equal(2))
	NewWithT(t).Expect(gs[0].Identifier()).To(Equal("a"))
	NewWithT(t).Expect(gs[1].Identifier()).To(Equal("b"))
}
