//go:generate go run ./std/main.go
package dumper

import (
	"bufio"
	"bytes"
	_ "embed"
	"sync"
	"sync/atomic"

	"github.com/xoctopus/typex/pkgutil"
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

type ImportTracker interface {
	// Track adds package path and name
	Track(string)
	// Package returns the ref leader of package path
	Package(string) string
	// Range traverse imports
	Range(func(Import))
	// Init adjust imported package alias
	Init()
}

func NewImportTracker(self string) ImportTracker {
	i := &tracker{
		imports: make(map[string]*Import),
		names:   make(map[string][]*Import),
		self:    self,
	}
	return i
}

type tracker struct {
	imports     map[string]*Import
	names       map[string][]*Import
	self        string
	once        sync.Once
	initialized atomic.Bool
}

func (t *tracker) Track(path string) {
	must.BeTrueF(
		!t.initialized.Load(),
		"cannot add package path to tracker after initialization",
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
	t.names[p.Name()] = append(t.names[p.Name()], i)
}

func (t *tracker) Package(path string) string {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot fetch package reference before tracker initialization",
	)

	imp, ok := t.imports[path]
	must.BeTrueF(ok, "imported package %s not be tracked", path)
	if path == t.self {
		return ""
	}
	return imp.alias
}

func (t *tracker) Range(f func(Import)) {
	must.BeTrueF(
		t.initialized.Load(),
		"cannot range imports before tracker initialization",
	)
	for _, i := range t.imports {
		if i.path == t.self {
			continue
		}
		f(*i)
	}
}

func (t *tracker) Init() {
	t.once.Do(func() {
		for _, list := range t.names {
			if len(list) > 1 {
				for _, _ = range list {
				}
			}
		}
		t.initialized.Store(true)
	})
}

type Import struct {
	path  string
	name  string
	alias string
}

func (i *Import) Path() string {
	return i.path
}

func (i *Import) String() string {
	if i.alias == i.name {
		return "\"" + i.path + "\""
	}
	return i.alias + " \"" + i.path + "\""
}

func (i *Import) IsStd() bool {
	_, ok := stds[i.path]
	return ok
}
