package genx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/internal/dumper"
	"github.com/xoctopus/genx/pkg/snippet"
)

type GeneratorNewer interface {
	New(Context) Generator
}

type Generator interface {
	Identifier() string
	Generate(Context, types.Type) error
}

// AggregationGeneratorMarker marks the generator as aggregated.
// the generated files will be aggregated into a single file
type AggregationGeneratorMarker interface {
	aggregation()
}

// GlobalGeneratorMarker marks the generator as global regardless of annotations.
type GlobalGeneratorMarker interface {
	global()
}

type GenerateNewer interface {
	Newer(Context) Generator
}

type Executor interface {
	Execute(context.Context, ...Generator) error
}

type Context interface {
	IsZero() bool

	Context() context.Context

	Packages() *pkgx.Packages
	Package() pkgx.Package
	PackageByPath(string) pkgx.Package
	// PackageByPos(token.Pos) pkgx.Package

	Render(snippet.Snippet)
}

type Args struct {
	Entrypoint []string
	Workdir    string
}

func NewContext(args *Args) Executor {
	ctx := pkgx.CtxWorkdir.With(context.Background(), args.Workdir)

	return &genc{
		args: args,
		pkgs: pkgx.NewPackages(ctx, args.Entrypoint...),
	}
}

type genc struct {
	args *Args
	pkgs *pkgx.Packages
	gens map[string]struct{} // generated suffix
	curr pkgx.Package

	ctx  func() context.Context
	file *genf
}

func (x *genc) IsZero() bool {
	return x.ctx == nil || x.file == nil || x.curr == nil
}

func (x *genc) Packages() *pkgx.Packages {
	return x.pkgs
}

func (x *genc) Package() pkgx.Package {
	return x.curr
}

func (x *genc) PackageByPath(path string) pkgx.Package {
	return x.pkgs.Package(path)
}

func (x *genc) Execute(ctx context.Context, generators ...Generator) error {
	x.gens = make(map[string]struct{})
	for _, g := range generators {
		x.gens["_genx_"+g.Identifier()+".go"] = struct{}{}
	}

	for path := range x.pkgs.Directs {
		p := x.pkgs.Package(path)
		must.NotNilF(p, "package is not found: %s", path)
		if err := x.exec(ctx, p, generators...); err != nil {
			return err
		}
	}
	return nil
}

func (x *genc) exec(ctx context.Context, p pkgx.Package, gs ...Generator) error {
	tags := p.Doc().Tags()
	ignores := tags["genx:ignore"]

	generators := map[string]Generator{}
	for _, g := range gs {
		_, ok := generators[g.Identifier()]
		must.BeTrueF(!ok, "duplicated generator: '%s'", g.Identifier())
		generators[g.Identifier()] = g
	}

	// filter ignored
	// eg:
	//	genx:ignore=x,y,z generator x y and z will be skipped
	//	genx:enum=false enum will be skipped
	for name := range generators {
		if slices.Contains(ignores, name) {
			delete(generators, name)
		}
		for tag, values := range tags {
			if tag == "genx:"+name && slices.Contains(values, "false") {
				delete(generators, name)
			}
		}
	}

	// must be generated defined in package document
	globals := make(map[string]Generator)
	for _, tag := range p.Doc().TagKeys() {
		if strings.HasPrefix(tag, "genx:") {
			name := strings.TrimPrefix(tag, "genx:")
			if g, ok := generators[name]; ok {
				globals[name] = g
				delete(generators, name)
			}
		}
	}

	aggregated := make(map[string][]*genc)
	for _, group := range []struct {
		generators map[string]Generator
		global     bool
	}{
		{globals, true},
		{generators, false},
	} {
		for _, g := range group.generators {
			xp := &genc{
				args: x.args,
				pkgs: x.pkgs,
				curr: p,
				gens: x.gens,
			}
			cs, err := xp.genpkg(ctx, g, group.global)
			if err != nil {
				return err
			}
			if len(cs) > 0 {
				aggregated[g.Identifier()] = cs
			}
		}
	}

	for name, gcs := range aggregated {
		g := globals[name]
		if g == nil {
			g = generators[name]
		}
		if _, ok := g.(AggregationGeneratorMarker); ok {
			filename := "zz_" + stringsx.LowerSnakeCase(name) + "_genx_" + name + ".go"
			xf := newgenf(p, name, "")

			trackers := make([]dumper.ImportTracker, 0, len(gcs))
			snippets := make([]snippet.Snippet, 0)
			for _, c := range gcs {
				trackers = append(trackers, dumper.From(c.ctx()))
				snippets = append(snippets, c.file.snippets...)
			}
			xf.snippets = []snippet.Snippet{snippet.Snippets(snippet.Block("\n"), snippets...)}

			merged := dumper.With(ctx, dumper.MergeTrackers(trackers...))
			if err := xf.write(merged, filename); err != nil {
				return err
			}
			continue
		}

		for _, xc := range gcs {
			filename := stringsx.LowerSnakeCase(xc.file.typ) + "_genx_" + g.Identifier() + ".go"
			if err := xc.file.write(xc.ctx(), filename); err != nil {
				return err
			}
		}
	}

	return nil
}

func (x *genc) genpkg(ctx context.Context, g Generator, global bool) ([]*genc, error) {
	prefix := "genx:" + g.Identifier()
	generated := make([]*genc, 0)

	for t := range x.curr.TypeNames().Elements() {
		pos := t.Node().Pos()
		filename := x.curr.FileSet().File(pos).Position(pos).Filename

		// skip generated files
		skip := false
		for suffix := range x.gens {
			if strings.HasSuffix(filename, suffix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		tags := t.Doc().Tags()
		values, ok := tags[prefix]
		// marked disable clearly
		if ok && slices.Contains(values, "false") {
			continue
		}

		// if not global must mark enable clearly
		if !global && !ok {
			continue
		}

		xf := &genc{
			args: x.args,
			pkgs: x.pkgs,
			gens: x.gens,
			curr: x.curr,
			file: newgenf(x.curr, g.Identifier(), t.TypeName()),
			ctx: sync.OnceValue(func() context.Context {
				return dumper.WithEntry(ctx, x.curr.Path())
			}),
		}
		if err := xf.gen(x.New(g), t.Type()); err != nil {
			return nil, err
		}
		if xf.file.IsNil() {
			continue
		}
		generated = append(generated, xf)
	}

	return generated, nil
}

func (x *genc) gen(g Generator, t types.Type) error {
	return g.Generate(x, t)
}

func (x *genc) New(g Generator) Generator {
	if newer, ok := g.(GeneratorNewer); ok {
		return newer.New(x)
	}
	return reflect.New(reflectx.Indirect(reflect.ValueOf(g)).Type()).Interface().(Generator)
}

func (x *genc) Render(s snippet.Snippet) {
	x.file.render(s)
}

func (x *genc) Context() context.Context {
	if x.ctx != nil {
		return x.ctx()
	}
	return context.Background()
}

func newgenf(p pkgx.Package, generator, typename string) *genf {
	return &genf{
		name: generator,
		typ:  typename,
		pkg:  p,
	}
}

type genf struct {
	name     string
	typ      string
	pkg      pkgx.Package
	snippets []snippet.Snippet
}

func (x *genf) IsNil() bool {
	return len(x.snippets) == 0
}

func (x *genf) render(s snippet.Snippet) {
	x.snippets = append(x.snippets, s)
}

func (x *genf) write(ctx context.Context, filename string) error {
	body := bytes.NewBuffer(nil)

	for code := range snippet.Snippets(
		snippet.NewLine(1),
		snippet.Poster(x.pkg.Unwrap().Name(), "genx:"+x.name),
		snippet.Imports(ctx),
	).Fragments(ctx) {
		body.WriteString(code)
	}

	body.WriteRune('\n')

	for _, s := range x.snippets {
		for code := range s.Fragments(ctx) {
			body.WriteString(code)
		}
	}

	data := body.Bytes()
	text := bytes.Split(data, []byte("\n"))

	fileset := token.NewFileSet()
	f, err := parser.ParseFile(
		fileset,
		filename,
		data,
		parser.ParseComments|parser.SkipObjectResolution|parser.AllErrors,
	)

	var serr scanner.ErrorList
	if err != nil && errors.As(err, &serr) && serr.Len() > 0 {
		e := serr[0]
		line, column := e.Pos.Line, e.Pos.Column-1

		for i := line - 10; i < line; i++ {
			if i > 0 {
				_, _ = fmt.Printf("%4d:", i+1)
				if len(text[i]) > 0 {
					_, _ = fmt.Printf(" %s\n", text[i])
				} else {
					_, _ = fmt.Printf("\n")
				}
			}
		}
		fmt.Printf("      %s↑\n", strings.Repeat(" ", column))
		fmt.Println(e.Msg)
		return err
	}

	output := must.NoErrorV(os.OpenFile(
		filepath.Join(x.pkg.SourceDir(), filename),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644,
	))
	defer func() { _ = output.Close() }()

	return format.Node(output, fileset, f)
}
