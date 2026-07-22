package uintstr

import (
	"bytes"
	"cmp"
	_ "embed"
	"go/types"
	"log"
	"runtime/debug"

	"github.com/xoctopus/x/misc/timer"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

//go:embed uintstr.tpl
var template []byte

func init() {
	genx.Register(&g{})
}

type g struct {
	genx.AggregationGeneratorMarker
}

func (g) Identifier() string {
	return "uintstr"
}

func (x *g) Version() string {
	v := ""
	if i, ok := debug.ReadBuildInfo(); ok {
		v = i.Main.Version
	}
	return cmp.Or(v, "devel")
}

func (x *g) Generate(c genx.Context, t types.Type) error {
	cost := timer.Span()

	tn, ok := t.(*types.Named)
	if !ok {
		return nil
	}

	if tn.TypeParams().Len() > 0 {
		return nil
	}

	b, ok := tn.Obj().Type().Underlying().(*types.Basic)
	if !ok {
		return nil
	}

	ctx := c.Context()
	switch b.Kind() {
	case types.Uint64, types.Uint32, types.Uint16, types.Uint8, types.Uint:
		log.Printf("genx:%s %s\n", x.Identifier(), t.String())
		defer func() {
			log.Printf("==> cost: %fs", cost().Seconds())
		}()

		c.Render(s.Template(
			bytes.NewReader(template),
			s.Arg(ctx, "T", s.IdentTT(ctx, t)),
			s.ArgExposeUnsafe(ctx, "strconv", "FormatUint").WithName("strconv.FormatUint"),
			s.ArgExposeUnsafe(ctx, "strconv", "ParseUint").WithName("strconv.ParseUint"),
		))
		return nil
	default:
		return nil
	}
}
