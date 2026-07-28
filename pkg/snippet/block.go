package snippet

import (
	"context"
	"fmt"
	"iter"
	"strconv"
)

// Block creates a Snippet from a raw string block.
func Block(v string) Snippet {
	return block(v)
}

// BlockF creates a Snippet from a formatted string block.
func BlockF(v string, args ...any) Snippet {
	return Block(fmt.Sprintf(v, args...))
}

// BlockRaw creates a Snippet from a string, quoting it automatically.
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
