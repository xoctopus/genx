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
	"strings"
	"sync"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/internal/dumper"
	"github.com/xoctopus/genx/pkg/docx"
	"github.com/xoctopus/genx/pkg/snippet"
)

type GeneratorNewer interface {
	New(Context) Generator
}

type Generator interface {
	Identifier() string
	Generate(Context, types.Type) error
}

type Versioned interface {
	Version() string
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
	type generator struct {
		Generator
		global bool
	}

	var (
		doc        = docx.Parse(p.PackageDoc())
		directives = doc.Directives() // global directives
		generators = make(map[string]generator)
		aggregated = make(map[string][]*genc)
	)

	for _, g := range gs {
		_, ok := generators[g.Identifier()]
		must.BeTrueF(!ok, "duplicated generator: '%s'", g.Identifier())
		// trim global disabled
		gg := generator{Generator: g}
		if d, ok := directives[g.Identifier()]; ok && d.Enabled() {
			gg.global = true
		}
		generators[g.Identifier()] = gg
	}

	for _, g := range generators {
		xp := &genc{
			args: x.args,
			pkgs: x.pkgs,
			curr: p,
			gens: x.gens,
		}
		cs, err := xp.genpkg(ctx, g.Generator, g.global)
		if err != nil {
			return err
		}
		if len(cs) > 0 {
			aggregated[g.Identifier()] = cs
		}
	}

	for name, gcs := range aggregated {
		g := generators[name].Generator
		if _, ok := g.(AggregationGeneratorMarker); ok {
			filename := "zz" + "_genx_" + stringsx.LowerSnakeCase(name) + ".go"

			version := ""
			if v, ok := g.(Versioned); ok {
				version = v.Version()
			}
			xf := newgenf(p, name, version, "")

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

		directives := docx.Parse(t.Doc()).Directives()
		gg, ok := directives[g.Identifier()]
		if ok && !gg.Enabled() {
			continue
		}
		if !ok && !global {
			continue
		}

		version := ""
		if v, ok := g.(Versioned); ok {
			version = v.Version()
		}

		xf := &genc{
			args: x.args,
			pkgs: x.pkgs,
			gens: x.gens,
			curr: x.curr,
			file: newgenf(x.curr, g.Identifier(), version, t.TypeName()),
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

func newgenf(p pkgx.Package, generator, version, typename string) *genf {
	return &genf{
		name:    generator,
		version: version,
		typ:     typename,
		pkg:     p,
	}
}

type genf struct {
	name     string
	version  string
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
		snippet.Poster(x.pkg.Unwrap().Name(), "genx:"+x.name, x.version),
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

	if se, ok := errors.AsType[scanner.ErrorList](err); ok && se.Len() > 0 {
		e := (se)[0]
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
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644,
	))
	defer func() { _ = output.Close() }()

	return format.Node(output, fileset, f)
}
