package codex

import (
	"bytes"
	_ "embed"
	"go/types"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

//go:embed codex.go.tpl
var template []byte

func init() {
	genx.Register(&g{})
}

type g struct {
	errors *Errors
}

func (x *g) Identifier() string {
	return "errx_code"
}

func (x *g) New(c genx.Context) genx.Generator {
	return &g{errors: NewErrors(c)}
}

func (x *g) Generate(c genx.Context, t types.Type) error {
	if e, ok := x.errors.Resolve(t); ok {
		if e.IsValid() {
			x.generate(c, e)
			return nil
		}
	}
	return nil
}

func (x *g) generate(c genx.Context, e *Error) {
	ctx := c.Context()
	ident := s.IdentTT(ctx, e.typ)

	args := []*s.TArg{
		// @def CodeType
		s.Arg(ctx, "CodeType", ident),
		// @def fmt.Sprintf
		s.ArgExpose(ctx, "fmt", "Sprintf"),
		// @def CodeMessageCases
		s.Arg(ctx, "CodeMessageCases", e.CodeMessageCases(ctx)),
		// @def UnknownCodeFormat
		s.Arg(ctx, "UnknownCodeFormat", s.BlockRaw("["+e.def+":%d] unknown")),
	}

	c.Render(s.Template(bytes.NewReader(template), args...))
}
