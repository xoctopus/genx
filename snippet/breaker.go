package snippet

import (
	"context"
	"iter"
	"strings"
)

func Indent(n int) Snippet {
	return &breaker{repeats: n, sep: '\t'}
}

func NewLine(n int) Snippet {
	return &breaker{repeats: n, sep: '\n'}
}

type breaker struct {
	repeats int
	sep     rune
}

func (b *breaker) IsNil() bool {
	return b.repeats == 0
}

func (b *breaker) Fragments(_ context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		x := strings.Repeat(string(b.sep), b.repeats)
		if !yield(x) {
			return
		}
	}
}
