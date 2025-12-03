package dumper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"
	gopkg "golang.org/x/tools/go/packages"
)

var gPackages = syncx.NewXmap[string, *gopkg.Package]()

func Load(ctx context.Context, path string) (pkg *gopkg.Package) {
	defer func(path string) {
		if pkg != nil {
			gPackages.Store(path, pkg)
		}
	}(path)

	_path := path
	cfg := pkgx.Config(ctx)
	if strings.HasSuffix(path, "_test") {
		path = strings.TrimSuffix(_path, "_test")
		cfg.Tests = true
	}
	pkgs, err := gopkg.Load(cfg, path)
	must.NoErrorF(err, "failed to load %s", path)
	must.BeTrueF(len(pkgs) > 0, "no packages loaded")
	for i := range pkgs {
		if pkgs[i].PkgPath == _path {
			pkg = pkgs[i]
			must.BeTrueF(len(pkg.Errors) == 0, "failed to load package %s", _path)
			break
		}
	}
	must.BeTrueF(pkg != nil, "failed to load package %s", _path)
	return pkg
}

func NewImport(ctx context.Context, path string) *Import {
	pkg := (*gopkg.Package)(nil)
	if p, ok := gPackages.Load(path); ok {
		pkg = p
	} else {
		pkg = Load(ctx, path)
	}
	i := &Import{
		pkg:   pkg,
		path:  path,
		name:  pkg.Name,
		alias: pkg.Name,
	}
	r := make([]rune, 0, len(path))

	// see https://go.dev/ref/spec#Import_declarations
	// go help importpath
	// an import path just contains alphabets and . _ - /
	for _, c := range []rune(path) {
		switch c {
		case '_', '.', '-', '/':
			r = append(r, '_')
		default:
			r = append(r, c)
		}
	}
	path = string(r)
	for _, s := range strings.Split(path, "_") {
		if s != "" {
			i.seps = append(i.seps, s)
		}
	}
	i.cut = 1

	return i
}

type Import struct {
	pkg   *gopkg.Package
	path  string
	name  string
	alias string
	seps  []string
	cut   int
}

func (i *Import) Path() string {
	return i.path
}

func (i *Import) Alias() string {
	return i.alias
}

func (i *Import) Name() string {
	return i.name
}

func (i *Import) Code() string {
	if i.alias == i.name {
		return strconv.Quote(i.path)
	}
	return fmt.Sprintf("%s %q", i.alias, i.path)
}

func (i *Import) MakeAlias() {
	i.cut++
	cut0 := len(i.seps) - i.cut
	if cut0 < 0 {
		cut0 = 0
	}
	i.alias = strings.Join(i.seps[cut0:], "_")
}
