package snippet

import (
	"context"
	"iter"
)

type Snippet interface {
	IsNil() bool
	Fragments(ctx context.Context) iter.Seq[string]
}

type Tracker interface {
	Track() string
}

func Snippets(sep Snippet, ss ...Snippet) Snippet {
	return &snippets{
		sep: sep,
		ss:  ss,
	}
}

func Compose(ss ...Snippet) Snippet {
	return &snippets{
		ss: ss,
	}
}

type snippets struct {
	sep Snippet
	ss  []Snippet
}

func (ss *snippets) IsNil() bool {
	return len(ss.ss) == 0
}

func (ss *snippets) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, si := range ss.ss {
			if si == nil || si.IsNil() {
				continue
			}
			if ss.sep != nil && !ss.sep.IsNil() && i > 0 {
				for s := range ss.sep.Fragments(ctx) {
					if !yield(s) {
						return
					}
				}
			}
			for s := range si.Fragments(ctx) {
				if !yield(s) {
					return
				}
			}
		}
	}
}

// type Writer interface {
// 	Render(Snippet)
// }
//
// func NewWriter(w io.Writer) Writer {
// 	return &writer{
// 		writer:  w,
// 		tracker: dumper.NewImportTracker(),
// 	}
// }
//
// type writer struct {
// 	writer  io.Writer
// 	tracker dumper.ImportTracker
// }
//
// func (w *writer) Render(s Snippet) {
// 	if s == nil || s.IsNil() {
// 		return
// 	}
//
// 	if !s.IsNil() {
// 		for code := range s.Fragments(dumper.Context.With(context.Background(), w.tracker)) {
// 			_, _ = io.WriteString(w.writer, code)
// 		}
// 	}
// }
