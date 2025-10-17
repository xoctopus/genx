package snippet

import (
	"context"
	"iter"
	"strings"
)

func Indent(n int) Snippet {
	return &separator{repeats: n, sep: '\t'}
}

func NewLine(n int) Snippet {
	return &separator{repeats: n, sep: '\n'}
}

type separator struct {
	repeats int
	sep     rune
}

func (b *separator) IsNil() bool {
	return b.repeats == 0
}

func (b *separator) Fragments(_ context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		x := strings.Repeat(string(b.sep), b.repeats)
		yield(x)
	}
}
