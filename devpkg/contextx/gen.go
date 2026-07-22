package contextx

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

//go:embed contextx.tpl
var template []byte

func init() {
	genx.Register(&g{})
}

type g struct {
	genx.AggregationGeneratorMarker
}

func (g) Identifier() string {
	return "ctx"
}

func (x *g) Version() string {
	v := ""
	if i, ok := debug.ReadBuildInfo(); ok {
		v = i.Main.Version
	}
	return cmp.Or(v, "devel")
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
	// generic context is not supported
	if t.TypeParams().Len() > 0 {
		return nil, true
	}

	ctx := c.Context()

	ident := s.IdentTT(ctx, t)
	pkgid := "github.com/xoctopus/x/contextx"

	// if tparams := t.TypeParams(); tparams.Len() > 0 {
	// 	e := c.Package().TypeNames().ElementByName(t.Obj().Name())
	// 	must.BeTrue(e != nil)
	// 	d := docx.Parse(e.Doc())
	// 	annotations := make(map[string]string)
	// 	for key, anno := range d.Annotations() {
	// 		if strings.HasPrefix(key, "arg") {
	// 			annotations[key] = anno[0].Text()
	// 		}
	// 	}
	// 	must.BeTrueF(len(annotations) == tparams.Len(), "mismatch type arguments count")
	// 	keys := slices.Collect(maps.Keys(annotations))
	// 	slices.Sort(keys)
	// 	for i, key := range keys {
	// 		argi := annotations[key]
	// 		must.BeTrueF(len(argi) > 0, "missing type argument %d", i)
	// 		targs = append(targs, argi)
	// 	}
	// }

	c.Render(s.Snippets(
		s.NewLine(1),
		s.Template(
			bytes.NewReader(template),
			s.ArgExposeUnsafe(ctx, pkgid, "From").WithName("contextx.From"),
			s.ArgExposeUnsafe(ctx, pkgid, "With").WithName("contextx.With"),
			s.ArgExposeUnsafe(ctx, pkgid, "Must").WithName("contextx.Must"),
			s.ArgExposeUnsafe(ctx, pkgid, "Carry").WithName("contextx.Carry"),
			s.Arg(ctx, "T", ident),
		),
	))
	return nil, false
}
