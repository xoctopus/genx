package snippet

import (
	"context"
	"iter"
)

func IsNil(s Snippet) bool {
	return s == nil || s.IsNil()
}

type Snippet interface {
	IsNil() bool
	Fragments(ctx context.Context) iter.Seq[string]
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

func Strings(tail string, sep string, list ...string) Snippet {
	if len(list) == 0 {
		return &Placeholder{}
	}
	ss := make([]Snippet, 0, len(list))
	for _, v := range list {
		ss = append(ss, Compose(BlockRaw(v), Block(tail)))
	}
	return Snippets(Block(sep), ss...)
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
			if IsNil(si) {
				continue
			}
			if !IsNil(ss.sep) && i > 0 {
				for s := range ss.sep.Fragments(ctx) {
					yield(s)
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

// Placeholder must be not nil and yield an empty string
type Placeholder struct{}

func (Placeholder) IsNil() bool {
	return false
}

func (Placeholder) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		yield("")
	}
}
