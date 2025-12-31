package docx

import (
	"bytes"
	_ "embed"
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/misc/must"

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

func (x *g) Generate(c genx.Context, t types.Type) error {
	switch t.Underlying().(type) {
	case nil, *types.Interface, *types.Alias, *types.Union:
		return nil // skipped
	}
	n, ok := t.(*types.Named)
	must.BeTrueF(ok, "accept *types.Named only but got %T", t)

	return x.generate(c, n)
}

func (x *g) generate(c genx.Context, t *types.Named) error {
	var (
		u     = t.Underlying()
		ctx   = c.Context()
		name  = t.Obj().Name()
		ident *s.TArg
	)

	if params := t.TypeParams(); params.Len() == 0 {
		ident = s.Arg(ctx, "T", s.Block(name))
	} else {
		names := make([]string, 0, params.Len())
		for i := range params.Len() {
			names = append(names, params.At(i).Obj().Name())
		}
		ident = s.Arg(ctx, "T", s.BlockF("%s[%s]", name, strings.Join(names, ", ")))
	}

	args := []*s.TArg{
		ident,
		s.Arg(ctx, "TDoc", x.docNamed(c, name)),
	}

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
				s.Arg(ctx, "FieldDocCases", x.docFields(c, y)),
				s.Arg(ctx, "AnonymousDoc", x.docAnonymous(c, y)),
			)...,
		)
	}

	c.Render(ss)
	return nil
}

func (x *g) doc(d *pkgx.Doc) string {
	if d == nil || len(d.Desc()) == 0 {
		return ""
	}
	lines := make([]string, 0)
	for _, desc := range d.Desc() {
		lines = append(lines, strconv.Quote(desc))
	}
	return strings.Join(lines, ", ")
}

func (x *g) docNamed(c genx.Context, typename string) s.Snippet {
	o := c.Package().TypeNames().ElementByName(typename)
	must.BeTrueF(o != nil, "type '%s' not found in package $s", typename, c.Package().Path())
	return s.Block(x.doc(o.Doc()))
}

func (x *g) docFields(c genx.Context, p *types.Struct) s.Snippet {
	ss := make([]s.Snippet, 0, p.NumFields())

	for i := 0; i < p.NumFields(); i++ {
		f := p.Field(i)
		if !ast.IsExported(f.Name()) {
			continue
		}

		// if embedded field is not defined in current package, treated as a Field
		if f.Embedded() {
			n, ok := f.Type().(*types.Named)
			if ok && n.Obj().Pkg().Path() == c.Package().Path() {
				continue
			}
		}
		d := c.Package().DocOf(f.Pos())
		if _, ok := f.Type().(*types.Struct); ok {
			continue // skip inline struct
		}
		if _s, ok := f.Type().Underlying().(*types.Struct); ok && _s.NumFields() == 0 {
			continue // skip empty struct
		}

		ss = append(
			ss,
			s.BlockF(`case %q:
	return []string{%s}, true`, f.Name(), x.doc(d)),
		)
	}
	return s.Snippets(s.Block("\n"), ss...)
}

func (x *g) docAnonymous(c genx.Context, p *types.Struct) s.Snippet {
	ss := make([]s.Snippet, 0, p.NumFields())
	ctx := c.Context()
	for i := 0; i < p.NumFields(); i++ {
		f := p.Field(i)
		if !f.Anonymous() {
			continue
		}
		if _s, ok := f.Type().Underlying().(*types.Struct); ok && !hasExported(_s) {
			continue
		}
		prefix := ""
		if d := c.Package().DocOf(f.Pos()); d != nil && len(d.Desc()) > 0 {
			prefix = d.Desc()[0]
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
