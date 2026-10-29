package codex

import (
	"bytes"
	_ "embed"
	"go/types"
	"log"

	"github.com/xoctopus/x/misc/timer"

	"github.com/xoctopus/genx/devpkg/helper"
	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

var (
	//go:embed codex.message.tpl
	template []byte
	//go:embed codex.status-describer.tpl
	templateDescriber []byte
	//go:embed codex.identifier.tpl
	templateIdentifier []byte
	//go:embed codex.hint.tpl
	templateHint []byte
	//go:embed codex.values.tpl
	templateValues []byte

	identifier = "code"
)

func init() {
	genx.Register(&g{})
}

type g struct {
	errors *Errors
}

func (x *g) Identifier() string {
	return identifier
}

func (x *g) Version() string {
	return helper.VersionFor("github.com/xoctopus/genx")
}

func (x *g) New(c genx.Context) genx.Generator {
	return &g{errors: NewErrors(c)}
}

func (x *g) Generate(c genx.Context, t types.Type) error {
	if e, ok := x.errors.Resolve(t); ok {
		if e.IsValid() {
			cost := timer.Span()
			log.Printf("genx:%s %s", x.Identifier(), t)
			x.generate(c, e)
			log.Printf("==> cost: %fs", cost().Seconds())
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
		s.ArgExposeUnsafe(ctx, "fmt", "Sprintf"),
		// @def CodeMessageCases
		s.Arg(ctx, "CodeMessageCases", e.CodeMessageCases(ctx)),
		// @def UnknownCodeFormat
		s.Arg(ctx, "UnknownCodeFormat", s.BlockRaw("["+e.domain+":%d] unknown")),
	}

	ss := []s.Snippet{
		s.Template(bytes.NewReader(template), args...),
	}

	if len(e.statuses) > 0 {
		args = append(
			args,
			// @def StatusDescriberCases
			s.Arg(ctx, "StatusDescriberCases", e.StatusDescriberCases(ctx)),
		)
		ss = append(ss, s.Template(bytes.NewReader(templateDescriber), args...))
	}

	args = append(
		args,
		// @def CodeIdentifierCases
		s.Arg(ctx, "CodeIdentifierCases", e.CodeIdentifierCases(ctx)),
	)
	ss = append(ss, s.Template(bytes.NewReader(templateIdentifier), args...))

	if len(e.hints) > 0 {
		args = append(
			args,
			// @def CodeHintCases
			s.Arg(ctx, "CodeHintCases", e.CodeHintCases(ctx)),
		)
		ss = append(ss, s.Template(bytes.NewReader(templateHint), args...))
	}

	args = append(
		args,
		// @def CodeValues
		s.Arg(ctx, "CodeValues", e.CodeValues(ctx)),
	)
	ss = append(ss, s.Template(bytes.NewReader(templateValues), args...))

	c.Render(s.Snippets(s.NewLine(1), ss...))
}
