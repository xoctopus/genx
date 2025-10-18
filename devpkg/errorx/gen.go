package errorx

import (
	"bytes"
	_ "embed"
	"go/types"

	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

//go:embed errorx.go.tpl
var template []byte

func init() {
	genx.Register(&g{})
}

type g struct {
	errors *Errors
}

func (x *g) Identifier() string {
	return "code_error"
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
		// @def MessagesVar
		s.Arg(ctx, "MessagesVar", s.Block(stringsx.LowerCamelCase(e.name)+"Messages")),
		// @def fmt.Sprintf
		s.ArgExpose(ctx, "fmt", "Sprintf"),
		// @def github.com/pkg/errors.As
		s.ArgExpose(ctx, "github.com/pkg/errors", "WithStack"),
		// @def github.com/pkg/errors.WithStack
		s.ArgExpose(ctx, "github.com/pkg/errors", "As"),
		// @def MessagesKeyValue
		s.Arg(ctx, "MessageKeyValues", e.MessageKeyValues(ctx)),
	}

	c.Render(s.Template(bytes.NewReader(template), args...))
}
