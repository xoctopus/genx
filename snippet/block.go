package snippet

import (
	"context"
	"fmt"
	"iter"
)

func Block(v string) Snippet {
	return block(v)
}

func BlockF(v string, args ...any) Snippet {
	return Block(fmt.Sprintf(v, args...))
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
