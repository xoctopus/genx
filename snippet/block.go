package snippet

import (
	"context"
	"fmt"
	"iter"
	"strconv"
)

func Block(v string) Snippet {
	return block(v)
}

func BlockF(v string, args ...any) Snippet {
	return Block(fmt.Sprintf(v, args...))
}

func BlockRaw(v string) Snippet {
	return block(strconv.Quote(v))
}

type block string

func (b block) IsNil() bool {
	return len(b) == 0
}

func (b block) Fragments(_ context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		yield(string(b))
	}
}
