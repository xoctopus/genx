package snippet

import (
	"context"
	"iter"
)

// IsNil checks if a Snippet is nil or its underlying implementation reports as nil.
func IsNil(s Snippet) bool {
	return s == nil || s.IsNil()
}

// Snippet represents a reusable code fragment generator.
type Snippet interface {
	// IsNil returns true if the snippet has no content to generate.
	IsNil() bool
	// Fragments yields a sequence of string fragments that make up the snippet.
	Fragments(ctx context.Context) iter.Seq[string]
}

// Snippets joins multiple Snippets together, separated by the given sep Snippet.
func Snippets(sep Snippet, ss ...Snippet) Snippet {
	return &snippets{
		sep: sep,
		ss:  ss,
	}
}

// Compose combines multiple Snippets sequentially without any separator.
func Compose(ss ...Snippet) Snippet {
	return &snippets{
		ss: ss,
	}
}

// Strings creates a Snippet from a list of strings.
// Each string is appended with the specified tail, and the resulting blocks
// are joined together using the given separator.
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
