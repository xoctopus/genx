package snippet

import (
	"context"
	"iter"
	"strings"

	"github.com/xoctopus/genx/internal/dumper"
)

func Imports(ctx context.Context) Snippet {
	tracker := dumper.TrackerFromContext(ctx)
	// when write import segment. invoke Init to finish tracking
	tracker.Init()

	s := &imports{}
	mod := tracker.Module()

	for i := range tracker.Range {
		if dumper.IsStd(i.Path()) {
			s.stds = append(s.stds, i)
			continue
		}
		if mod != "" && strings.HasPrefix(i.Path(), mod) {
			s.projects = append(s.projects, i)
			continue
		}
		s.generals = append(s.generals, i)
	}

	// cmp := func(x, y dumper.Import) int {
	// 	if x.Path() < y.Path() {
	// 		return -1
	// 	}
	// 	if x.Path() > y.Path() {
	// 		return 1
	// 	}
	// 	return 0
	// }

	// slices.SortFunc(s.stds, cmp)
	// slices.SortFunc(s.generals, cmp)
	// slices.SortFunc(s.projects, cmp)

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
		b := &strings.Builder{}
		grouped := [][]dumper.Import{i.stds, i.generals, i.projects}

		b.WriteString("import (\n")

		written := 0
		for _, group := range grouped {
			if len(group) == 0 {
				continue
			}
			if written > 0 {
				b.WriteString("\n")
			}
			for _, s := range group {
				b.WriteString("\t")
				b.WriteString(s.Code())
				b.WriteString("\n")
			}
			written++
		}
		b.WriteString(")\n\n")

		yield(b.String())
	}
}
