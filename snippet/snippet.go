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
