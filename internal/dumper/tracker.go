package dumper

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/misc/must"
)

var (
	//go:embed std.list
	std []byte

	stds = make(map[string]struct{})
)

func init() {
	scanner := bufio.NewScanner(bytes.NewBuffer(std))
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			stds[line] = struct{}{}
		}
	}
}

func IsStd(path string) bool {
	_, ok := stds[path]
	return ok
}

type ImportTracker interface {
	// Track adds package path and name
	Track(context.Context, string)
	// TrackCustom adds package tracking with custom alias
	TrackCustom(context.Context, string, string)
	// PackageName returns the ref leader of package path
	PackageName(string) string
	// Range traverse imports
	Range(func(Import) bool)
	// Init adjust imported package alias
	Init()
	// Entry returns tracker's entry package path
	Entry() string
	// Module return tracker's module path
	Module() string
}

func MergeTrackers(trackers ...ImportTracker) ImportTracker {
	for i, j := 0, 0; i < len(trackers); i++ {
		if trackers[i] != nil {
			trackers[j] = trackers[i]
			j++
		}
	}
	must.BeTrueF(len(trackers) > 0, "trackers must not be empty")

	entry := trackers[0].Entry()
	merged := NewImportTracker(entry)
	for _, t := range trackers {
		tr := t.(*tracker)
		must.BeTrueF(
			!tr.initialized.Load(),
			"cannot track package to tracker after initialization",
		)
		must.BeTrueF(
			t.Entry() == entry,
			"cannot merge trackers with different entry",
		)
		for _, i := range tr.imports {
			merged.Track(context.Background(), i.path)
		}
	}
	return merged
}

func NewImportTracker(entry string) ImportTracker {
	pkg := NewImport(context.Background(), entry).pkg
	i := &tracker{
		imports: make(map[string]*Import),
		names:   make(map[string][]*Import),
		entry:   pkg.PkgPath,
		module:  pkg.Module.Path,
	}
	return i
}

type tracker struct {
	imports     map[string]*Import
	names       map[string][]*Import
	entry       string
	module      string
	paths       []string // ordered: std, external, entry module
	once        sync.Once
	initialized atomic.Bool
}

var (
	_ ImportTracker = (*tracker)(nil)
	_ pkgx.PkgNamer = (*tracker)(nil)
)

func (t *tracker) Track(ctx context.Context, path string) {
	t.TrackCustom(ctx, path, "")
}

func (t *tracker) TrackCustom(ctx context.Context, path string, alias string) {
	must.BeTrueF(
		!t.initialized.Load(),
		"cannot track package to tracker after initialization",
	)
	if path == "" || path == t.entry {
		return
	}

	if _, ok := t.imports[path]; ok {
		return
	}

	i := NewImport(ctx, path)
	t.imports[path] = i
	if alias != "" {
		i.alias = alias
	}
	t.names[i.alias] = append(t.names[i.Name()], i)
}

func (t *tracker) PackageName(path string) string {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot fetch package reference before tracker initialization",
	)

	if path == t.entry {
		return ""
	}

	i, ok := t.imports[path]
	must.BeTrueF(ok, "imported package %s not be tracked", path)
	return i.alias
}

func (t *tracker) Range(f func(Import) bool) {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot range imports before tracker initialization",
	)

	for _, path := range t.paths {
		i := t.imports[path]
		f(*i)
	}
}

func (t *tracker) Init() {
	t.once.Do(func() {
		t.paths = make([]string, 0, len(t.imports))
		for _, i := range t.imports {
			t.paths = append(t.paths, i.path)
		}

		sort.Slice(t.paths, func(i, j int) bool {
			pi := t.imports[t.paths[i]].path
			pj := t.imports[t.paths[j]].path
			si, sj := IsStd(pi), IsStd(pj)

			if si == sj {
				return pi < pj
			}
			return si
		})

		check := func() string {
			for name := range t.names {
				if len(t.names[name]) > 1 {
					return name
				}
			}
			return ""
		}

		for name := check(); name != ""; name = check() {
			list := t.names[name]
			delete(t.names, name)
			for _, i := range list {
				if !IsStd(i.path) {
					i.MakeAlias()
				}
				t.names[i.Alias()] = append(t.names[i.Alias()], i)
			}
		}
		t.initialized.Store(true)
	})
}

func (t *tracker) Entry() string {
	return t.entry
}

func (t *tracker) Module() string {
	return t.module
}
