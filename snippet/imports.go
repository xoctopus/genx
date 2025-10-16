package snippet

import (
	"context"
	"iter"
	"slices"
	"strings"

	"github.com/xoctopus/genx/internal/dumper"
)

func Imports(ctx context.Context, mod string) Snippet {
	tracker := dumper.TrackerFromContext(ctx)
	// when write import segment. invoke Init to finish tracking
	tracker.Init()

	s := &imports{}

	tracker.Range(func(i dumper.Import) {
		if i.IsStd() {
			s.stds = append(s.stds, i)
			return
		}
		if mod != "" && strings.Contains(i.Path(), mod) {
			s.projects = append(s.projects, i)
			return
		}
		s.generals = append(s.generals, i)
	})

	cmp := func(x, y dumper.Import) int {
		if x.Path() < y.Path() {
			return -1
		}
		if x.Path() > y.Path() {
			return 1
		}
		return 0
	}

	slices.SortFunc(s.stds, cmp)
	slices.SortFunc(s.generals, cmp)
	slices.SortFunc(s.projects, cmp)

	return s
}

type imports struct {
	stds     []dumper.Import
	generals []dumper.Import
	projects []dumper.Import
}

func (i *imports) IsNil() bool {
	return len(i.stds)+len(i.generals)+len(i.projects) == 0
}

func (i *imports) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		grouped := [][]dumper.Import{i.stds, i.generals, i.projects}
		if !yield("import (\n") {
			return
		}

		written := 0
		for _, group := range grouped {
			if len(group) == 0 {
				continue
			}
			if written > 0 {
				if !yield("\n") {
					return
				}
			}
			for _, s := range group {
				if !yield("\t") {
					return
				}
				if !yield(s.String()) {
					return
				}
				if !yield("\n") {
					return
				}
			}
			written++
		}
		yield(")\n\n")
	}
}
