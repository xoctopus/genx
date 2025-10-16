package snippet

import (
	"context"
	"iter"
)

func Block(v string) Snippet {
	return block(v)
}

type block string

func (b block) IsNil() bool {
	return len(b) == 0
}

func (b block) Fragments(_ context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		if !yield(string(b)) {
			return
		}
	}
}
