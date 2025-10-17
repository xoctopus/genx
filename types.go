package genx

import (
	"bytes"
	"context"
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

	"github.com/pkg/errors"
	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/stringsx"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/genx/internal/dumper"
	"github.com/xoctopus/genx/snippet"
)

type GeneratorNewer interface {
	New(Context) Generator
}

type Generator interface {
	Identifier() string
	Generate(Context, types.Type) error
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

	Package() pkgx.Package
	PackageByPath(string) pkgx.Package
	// PackageByPos(token.Pos) pkgx.Package

	Render(snippet.Snippet)
}

type Args struct {
	Entrypoint []string
}

func NewContext(args *Args) Executor {
	return &genc{
		args: args,
		pkgs: pkgx.NewPackages(args.Entrypoint...),
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

	for path := range x.pkgs.Directs() {
		p := x.pkgs.Package(path)
		must.NotNilF(p, "package is not found: %s", path)
		if err := x.exec(ctx, p, generators...); err != nil {
			return err
		}
	}
	return nil
}

func (x *genc) exec(ctx context.Context, p pkgx.Package, generators ...Generator) error {
	tags := p.Doc().Tags()
	ignores := tags["genx:ignore"]

	for _, g := range generators {
		// eg: the following generator will be skipped
		// genx:enum=false
		// genx:ignore=enum

		if slices.Contains(ignores, g.Identifier()) {
			continue
		}

		skip := false
		for _, v := range tags["genx:"+g.Identifier()] {
			if v == "false" {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		xp := &genc{
			args: x.args,
			pkgs: x.pkgs,
			curr: p,
			gens: x.gens,
		}
		if err := xp.genpkg(ctx, g); err != nil {
			return err
		}
	}

	return nil
}

func (x *genc) genpkg(ctx context.Context, g Generator) error {
	prefix := "genx:" + g.Identifier()
	generated := syncx.NewSmap[string, *genc]()

	for t := range x.curr.TypeNames().Elements() {
		pos := t.Node().Pos()
		filename := x.curr.FileSet().File(pos).Position(pos).Filename
		for suffix := range x.gens {
			if strings.HasSuffix(filename, suffix) {
				continue
			}
		}

		tags := t.Doc().Tags()
		values, ok := tags[prefix]
		if !ok || slices.Contains(values, "false") {
			continue
		}

		xf := &genc{
			args: x.args,
			pkgs: x.pkgs,
			gens: x.gens,
			curr: x.curr,
			file: newgenf(x.curr, g.Identifier()),
			ctx: sync.OnceValue(func() context.Context {
				return dumper.WithTrackerContext(ctx, x.curr.Unwrap().Path(), x.curr.GoModule().Path)
			}),
		}
		if err := xf.gen(x.New(g), t.Type()); err != nil {
			return err
		}
		if xf.file.IsNil() {
			continue
		}
		generated.Store(
			stringsx.LowerSnakeCase(t.TypeName())+"_genx_"+g.Identifier()+".go",
			xf,
		)
	}

	for filename, xf := range generated.Range {
		if err := xf.file.write(xf.ctx(), filename); err != nil {
			return err
		}
	}
	return nil
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

func newgenf(p pkgx.Package, name string) *genf {
	return &genf{
		name: name,
		pkg:  p,
	}
}

type genf struct {
	name     string
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
		snippet.Poster(x.pkg.Unwrap().Name(), x.name),
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

		b := &strings.Builder{}

		for i := line - 10; i < line; i++ {
			if i > 0 {
				_, _ = fmt.Fprintf(b, "%4d:", i+1)
				if len(text[i]) > 0 {
					_, _ = fmt.Fprintf(b, " %s\n", text[i])
				} else {
					_, _ = fmt.Fprintf(b, "\n")
				}
			}
		}
		if column < 0 {
			column = 0
		}
		_, _ = fmt.Fprintf(b, "      %s↑\n", strings.Repeat(" ", column))
		_, _ = fmt.Fprintln(b, e.Msg)
		fmt.Print(b.String())
		return err
	}

	filename = filepath.Join(x.pkg.SourceDir(), filename)
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	return format.Node(output, fileset, f)
}
