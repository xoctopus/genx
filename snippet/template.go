package snippet

import (
	"bufio"
	"context"
	"io"
	"iter"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

func Arg(ctx context.Context, name string, s Snippet) *TArg {
	return &TArg{
		name:    name,
		snippet: s,
	}
}

func ArgFor[T any](ctx context.Context, name string) *TArg {
	return &TArg{
		name:    name,
		snippet: IdentFor[T](ctx),
	}
}

func ArgT[T any](ctx context.Context) *TArg {
	id := IdentFor[T](ctx).(*ident)
	return &TArg{
		name:    id.t.String(),
		snippet: id,
	}
}

func ArgExpose(ctx context.Context, path string, name string, targs ...Snippet) *TArg {
	return &TArg{
		name:    path + "." + name,
		snippet: Expose(ctx, path, name, targs...),
	}
}

type TArg struct {
	name    string
	snippet Snippet
}

func (a *TArg) WithName(name string) *TArg {
	a2 := *a
	a2.name = name
	return &a2
}

type segment struct {
	name   string
	text   []string
	args   map[string]*TArg
	offset int
}

func (s *segment) IsNil() bool {
	return s == nil || len(s.text) == 0
}

func (s *segment) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		lineno := -1
		for _, line := range s.text {
			lineno++
			var (
				raw     = false // string quoted by "
				quoted  = false // macro quoted by # as replacer
				macro   string
				runes   = []rune(line)
				newline = make([]rune, 0)
				column  = -1
				index   = -1
				whole   = false // flag replace whole line
			)
			for i := 0; i < len(runes); i++ {
				column++
				c := runes[i]
				switch c {
				case '"':
					raw = !raw // quoted text won't be replaced
				case '#':
					if raw {
						continue
					}
					quoted = !quoted
					if quoted {
						must.BeTrue(index == -1)
						index = i
						continue
					}
					must.BeTrue(index >= 0)
					macro = string(runes[index+1 : i])
					if strings.TrimSpace(line) == "#"+macro+"#" {
						whole = true
					}
					goto FinishMacro
				}
				if !quoted {
					newline = append(newline, c)
				}
				continue
			FinishMacro:
				if whole {
					newline = newline[:0]
				}

				arg, ok := s.args[macro]
				must.BeTrueF(
					ok && arg != nil && !arg.snippet.IsNil(),
					"template argument %s not found or nil at line:%d:col",
					macro, lineno, column,
				)
				for code := range arg.snippet.Fragments(ctx) {
					newline = append(newline, []rune(code)...)
				}
				index = -1
				if whole {
					break
				}
			}
			must.BeTrueF(index == -1, "unfinished replacer at line:%d", lineno)
			if newline[len(newline)-1] != '\n' {
				newline = append(newline, '\n')
			}
			if !yield(string(newline)) {
				return
			}
		}
	}
}

func Template(r io.Reader, args ...*TArg) Snippet {
	tpl := &template{args: make(map[string]*TArg)}

	for _, arg := range args {
		if arg != nil {
			tpl.args[arg.name] = arg
		}
	}
	defs := make(map[string]int)

	var (
		scanner = bufio.NewScanner(r)
		seg     *segment
		lineno  = 0
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineno++

		// skip empty line
		if len(line) == 0 {
			continue
		}

		// collect replacer
		if strings.HasPrefix(line, "@def ") {
			defs[strings.TrimSpace(line[5:])] = lineno
			continue
		}

		if strings.HasPrefix(line, "--") {
			// handle prev
			if (seg == nil || seg.IsNil()) && len(tpl.segments) > 0 {
				tpl.segments = tpl.segments[:len(tpl.segments)-1]
			}
			seg = &segment{offset: lineno}
			seg.name = strings.TrimSpace(line[2:])
			tpl.segments = append(tpl.segments, seg)
			continue
		}
		must.BeTrueF(
			seg != nil, "please define a segment in a newline "+
				"prefixed with `--` to entry a new segment at line %d",
			lineno,
		)
		seg.text = append(seg.text, line)
	}

	for def := range defs {
		_, ok := tpl.args[def]
		must.BeTrueF(
			ok,
			"template argument name must be given. %s not in arguments",
			def,
		)
	}

	return tpl
}

type template struct {
	args     map[string]*TArg
	segments []*segment
}

func (t *template) IsNil() bool {
	return t == nil || len(t.segments) == 0
}

func (t *template) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, s := range t.segments {
			s.args = t.args
			for line := range s.Fragments(ctx) {
				if !yield(line) {
					return
				}
			}
			if !yield("\n") {
				return
			}
		}
	}
}
