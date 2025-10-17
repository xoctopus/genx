package dumper

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/xoctopus/typex/namer"
	"github.com/xoctopus/typex/pkgutil"
	"github.com/xoctopus/x/misc/must"
	"golang.org/x/exp/maps"
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
	Track(string)
	// Package returns the ref leader of package path
	Package(string) string
	// Range traverse imports
	Range(func(Import) bool)
	// Init adjust imported package alias
	Init()
	// Entry returns tracker's entry package path
	Entry() string
	// Module return tracker's module path
	Module() string
}

func NewImportTracker(entry string, module string) ImportTracker {
	i := &tracker{
		imports: make(map[string]*Import),
		names:   make(map[string][]*Import),
		entry:   entry,
		module:  module,
	}
	return i
}

type tracker struct {
	imports     map[string]*Import
	names       map[string][]*Import
	entry       string
	module      string
	once        sync.Once
	initialized atomic.Bool
}

var (
	_ ImportTracker      = (*tracker)(nil)
	_ namer.PackageNamer = (*tracker)(nil)
)

func (t *tracker) Track(path string) {
	must.BeTrueF(
		!t.initialized.Load(),
		"cannot track package to tracker after initialization",
	)
	if path == "" {
		return
	}

	if _, ok := t.imports[path]; ok {
		return
	}

	p := pkgutil.New(path)
	i := &Import{
		path:  path,
		name:  p.Name(),
		alias: p.Name(),
	}
	t.imports[path] = i
	t.names[i.alias] = append(t.names[p.Name()], i)
}

func (t *tracker) Package(path string) string {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot fetch package reference before tracker initialization",
	)

	i, ok := t.imports[path]
	must.BeTrueF(ok, "imported package %s not be tracked", path)
	if path == t.entry {
		return ""
	}
	return i.alias
}

func (t *tracker) Range(f func(Import) bool) {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot range imports before tracker initialization",
	)
	paths := maps.Keys(t.imports)
	sort.Slice(paths, func(i, j int) bool {
		pi := t.imports[paths[i]].path
		pj := t.imports[paths[j]].path
		if IsStd(pi) {
			if IsStd(pj) {
				return pi < pj
			}
			return true
		}
		return !IsStd(pj) && pi < pj
	})

	for _, path := range paths {
		i := t.imports[path]
		if i.path == t.entry {
			continue
		}
		f(*i)
	}
}

func (t *tracker) Init() {
	t.once.Do(func() {
		for name, list := range t.names {
			if len(list) <= 1 {
				continue
			}

			t.names[name] = nil
			sort.Slice(list, func(i, j int) bool {
				return list[i].path < list[j].path
			})

			externals := make([]*Import, 0, len(list))
			for _, p := range list {
				if IsStd(p.path) {
					t.names[name] = append(t.names[name], p)
					continue
				}
				externals = append(externals, p)
			}

			for _, p := range list {
				parts := strings.Split(p.path, "/")
				for i := range len(parts) {
					alias := strings.Join(parts[len(parts)-i-1:], "_")
					alias = strings.Replace(alias, ".", "_", -1)
					if _, ok := t.names[alias]; !ok {
						p.alias = alias
						break
					}
				}
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

type Import struct {
	path  string
	name  string
	alias string
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
