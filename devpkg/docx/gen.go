package docx

import (
	"bytes"
	"cmp"
	_ "embed"
	"fmt"
	"go/types"
	"log"
	"runtime/debug"
	"strings"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/misc/timer"
	"github.com/xoctopus/x/slicex"

	"github.com/xoctopus/genx/pkg/docx"
	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

var (
	//go:embed tpls/docx.named.tpl
	tplNamed []byte
	//go:embed tpls/docx.struct.tpl
	tplStruct []byte
	//go:embed tpls/docx.embed.tpl
	tplEmbed []byte
)

func init() {
	genx.Register(&g{})
}

type g struct {
	genx.AggregationGeneratorMarker
}

func (x *g) Identifier() string {
	return "doc"
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
	log.Printf("genx:%s %s", x.Identifier(), t)

	skip := false
	defer func() {
		if skip {
			log.Printf("==> cost: %fs skipped", cost().Seconds())
		} else {
			log.Printf("==> cost: %fs", cost().Seconds())
		}
	}()

	switch t.(type) {
	case nil, *types.Interface, *types.Alias, *types.Union:
		skip = true
		return
	}
	switch t.Underlying().(type) {
	case nil, *types.Interface, *types.Alias, *types.Union:
		skip = true
		return
	}

	return x.generate(c, t.(*types.Named))
}

func (x *g) generate(c genx.Context, t *types.Named) error {
	var (
		u        = t.Underlying()
		ctx      = c.Context()
		typename = t.Obj().Name()
		ident    *s.TArg
	)

	if params := t.TypeParams(); params.Len() == 0 {
		ident = s.Arg(ctx, "T", s.Block(typename))
	} else {
		pnames := make([]string, 0, params.Len())
		for param := range params.TypeParams() {
			pnames = append(pnames, param.Obj().Name())
		}
		ident = s.Arg(ctx, "T", s.BlockF("%s[%s]", typename, strings.Join(pnames, ", ")))
	}

	args := []*s.TArg{ident, s.Arg(ctx, "TDoc", x.docNamed(c, typename))}

	var ss s.Snippet
	if y, ok := u.(*types.Struct); !ok || !hasExported(y) {
		ss = s.Template(
			bytes.NewReader(tplNamed),
			args...,
		)
	} else {
		ss = s.Template(
			bytes.NewReader(tplStruct),
			append(
				args,
				s.Arg(ctx, "FieldDocCases", x.docNamedFields(c, typename, y)),
				s.Arg(ctx, "AnonymousDoc", x.docEmbeddedFields(c, typename, y)),
			)...,
		)
	}

	c.Render(ss)
	return nil
}

func (x *g) doc(prefix string, d *docx.Meta) string {
	doc := fmt.Sprintf("%q", "")

	if d != nil {
		lines := slicex.Mapping(d.Descriptions(prefix), func(line string) string {
			return fmt.Sprintf("%q", line)
		})
		doc = strings.Join(lines, ",")
	}
	return doc
}

func (x *g) docNamed(c genx.Context, typename string) s.Snippet {
	o := c.Package().TypeNames().ElementByName(typename)
	must.NotNilF(o, "type '%s' not found in package $s", typename, c.Package().Path())
	return s.Block(x.doc(typename, docx.Parse(o.Doc())))
}

func (x *g) docNamedFields(c genx.Context, typename string, p *types.Struct) s.Snippet {
	ss := make([]s.Snippet, 0, p.NumFields())

	for f := range p.Fields() {
		// if embedded field is not defined in current package, treated as a Field
		if f.Embedded() {
			continue
		}

		if d := c.Package().FieldDoc(typename, f.Name()); d != nil {
			ss = append(
				ss,
				s.BlockF(`case %q:`, f.Name()),
				s.Compose(
					s.Indent(1),
					s.BlockF(`return []string{%s}, true`, x.doc(f.Name(), docx.Parse(d))),
				),
			)
		}
	}
	return s.Snippets(s.Block("\n"), ss...)
}

func (x *g) docEmbeddedFields(c genx.Context, typename string, p *types.Struct) s.Snippet {
	ss := make([]s.Snippet, 0, p.NumFields())
	ctx := c.Context()
	for f := range p.Fields() {
		if !f.Embedded() {
			continue
		}

		if _s, ok := f.Type().Underlying().(*types.Struct); ok && !hasExported(_s) {
			continue
		}

		prefix := ""
		if d := c.Package().FieldDoc(typename, f.Name()); d != nil {
			prefix = docx.Parse(d).Title(f.Name())
		}

		ref := "&v"
		if _, ok := f.Type().(*types.Pointer); ok {
			ref = "v"
		}
		ss = append(
			ss,
			s.Template(
				bytes.NewReader(tplEmbed),
				s.Arg(ctx, "Ref", s.Block(ref)),
				s.Arg(ctx, "Field", s.Block(f.Name())),
				s.Arg(ctx, "DocOf", s.ExposeUnsafe(ctx, "github.com/xoctopus/x/docx", "Of")),
				s.Arg(ctx, "Prefix", s.BlockRaw(prefix)),
			),
		)
	}
	return s.Snippets(s.Block("\n"), ss...)
}

func hasExported(s *types.Struct) bool {
	for f := range s.Fields() {
		if f.Exported() {
			return true
		}
	}
	return false
}
