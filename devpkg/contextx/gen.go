package contextx

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
	//go:embed contextx.tpl
	templateNormal []byte
	//go:embed contextx.generic.tpl
	templateGeneric []byte
)

func init() {
	genx.Register(&g{})
}

type g struct {
	genx.AggregationGeneratorMarker
}

func (g) Identifier() string {
	return "context"
}

func (x *g) Version() string {
	return helper.VersionFor("github.com/xoctopus/genx")
}

func (x *g) Generate(c genx.Context, t types.Type) (err error) {
	cost := timer.Span()
	skip := false

	log.Printf("genx:%s %s", x.Identifier(), t.String())
	defer func() {
		if skip {
			log.Printf("==> cost: %fs skipped", cost().Seconds())
		} else {
			log.Printf("==> cost: %fs", cost().Seconds())
		}
	}()

	tn, ok := t.(*types.Named)
	if !ok {
		skip = true
		return
	}

	err, skip = x.generate(c, tn)
	return
}

func (x *g) generate(c genx.Context, t *types.Named) (error, bool) {
	ctx := c.Context()

	pkgid := "github.com/xoctopus/x/contextx"
	template := templateNormal

	args := []*s.TArg{
		s.ArgExposeUnsafe(ctx, pkgid, "From").WithName("contextx.From"),
		s.ArgExposeUnsafe(ctx, pkgid, "With").WithName("contextx.With"),
		s.ArgExposeUnsafe(ctx, pkgid, "Must").WithName("contextx.Must"),
		s.ArgExposeUnsafe(ctx, pkgid, "Carry").WithName("contextx.Carry"),
		s.ArgExpose(ctx, "context", "Context").WithName("context.Context"),
	}

	if n := t.TypeParams().Len(); n > 0 {
		template = templateGeneric
		tPkgPath := ""

		if p := t.Obj().Pkg(); p != nil {
			tPkgPath = p.Path()
		}

		params := make([]s.Snippet, 0, n)
		names := make([]s.Snippet, 0, n)
		for i := range n {
			ti := t.TypeParams().At(i)
			obj := ti.Obj()

			params = append(
				params,
				s.Snippets(s.Block(" "), s.Block(obj.Name()), s.IdentTT(ctx, ti.Constraint())),
			)

			names = append(names, s.Block(obj.Name()))
		}

		args = append(
			args,
			s.ArgExposeUnsafe(ctx, pkgid, "Carrier").WithName("contextx.Carrier"),
			s.Arg(ctx, "T", s.ExposeUnsafe(ctx, tPkgPath, t.Obj().Name())).WithName("T"),
			s.Arg(ctx, "TParams", s.Snippets(s.Block(", "), params...)),
			s.Arg(ctx, "TNames", s.Snippets(s.Block(", "), names...)),
		)
	} else {
		args = append(args, s.Arg(ctx, "T", s.IdentTT(ctx, t)))
	}

	c.Render(s.Snippets(
		s.NewLine(1),
		s.Template(bytes.NewReader(template), args...),
	))
	return nil, false
}
