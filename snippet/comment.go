package snippet

import (
	"context"
	"iter"
	"strings"
)

func Document(name string, lines ...string) Snippet {
	c := Comments(lines...).(*comment)
	name = strings.TrimSpace(name)
	if len(c.lines) > 0 && len(name) > 0 {
		c.lines[0] = name + " " + c.lines[0]
	}
	c.inline = false
	return c
}

func InlineComment(line string) Snippet {
	c := Comments(strings.Split(line, "\n")...).(*comment)
	if len(c.lines) > 0 {
		c.inline = true
	}
	return c
}

func Comments(lines ...string) Snippet {
	c := &comment{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimSpace(line)

		if strings.Contains(line, "\n") {
			cc := Comments(strings.Split(line, "\n")...).(*comment)
			c.lines = append(c.lines, cc.lines...)
			continue
		}
		if len(line) > 0 {
			c.lines = append(c.lines, line)
		}
	}
	return c
}

func Directive(directive string, args ...string) Snippet {
	final := []string{directive}
	for _, arg := range args {
		if a := strings.TrimSpace(arg); a != "" {
			final = append(final, a)
		}
	}

	return &comment{
		inline:    false,
		directive: true,
		lines:     []string{strings.Join(final, " ")},
	}
}

type comment struct {
	inline    bool
	directive bool
	lines     []string
}

func (c *comment) IsNil() bool {
	return len(c.lines) == 0
}

func (c *comment) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, line := range c.lines {
			if c.directive {
				line = "//go:" + line
			} else {
				line = "// " + line
			}
			if !yield(line) {
				return
			}
			if i < len(c.lines)-1 {
				if !yield("\n") {
					return
				}
			}
		}
	}
}
