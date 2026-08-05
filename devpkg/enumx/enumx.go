package enumx

import (
	"context"
	"fmt"
	"go/constant"
	"go/types"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/pkg/docx"
	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

type option struct {
	name  string
	enum  string
	text  string
	value *pkgx.Constant

	mapping map[string]string
	extends map[string]string
}

func (o *option) quotes() []string {
	quotes := []string{strconv.Quote(o.name)}
	if len(o.enum) > 0 {
		quotes = append(quotes, strconv.Quote(o.enum))
	}
	return quotes
}

type Enum struct {
	typ     types.Type
	key     string
	unknown *pkgx.Constant
	options []*option

	storage string

	mapping map[string]string
	extends map[string]bool
}

func (e *Enum) IsValid() bool {
	return e.unknown != nil || len(e.options) > 0
}

// add adds option
func (e *Enum) add(c *pkgx.Constant) {
	if name := c.Name(); name != "_" {
		prefix := stringsx.UpperSnakeCase(e.key)
		if name == prefix+"_UNKNOWN" {
			e.unknown = c
			return
		}

		before, suffix, found := strings.Cut(name, "__")
		if !found || before != prefix {
			return
		}

		doc := docx.Parse(c.Doc())

		o := &option{
			value:   c,
			name:    suffix,
			text:    doc.Title(c.Name()),
			mapping: make(map[string]string),
			extends: make(map[string]string),
		}

		e.options = append(e.options, o)

		if annotations, ok := doc.AnnotationsByName(identifier); ok {
			for _, anno := range annotations {
				if anno.Key() == "enum" {
					o.enum = anno.Value()
					continue
				}
				method, key, found := strings.Cut(anno.Key(), ".")
				if !found {
					continue
				}
				k := strings.ToLower(key)
				if k == "string" || k == "value" || key == "text" {
					continue
				}
				switch strings.ToLower(method) {
				case "map":
					if _, ok := e.mapping[key]; ok {
						o.mapping[o.name] = anno.Value()
					}
				case "ext":
					k = stringsx.UpperCamelCase(key)
					e.extends[k] = true
					o.extends[k] = anno.Value()
				default:
				}
			}
		}
	}
}

// Values generates code snippet of const value list
func (e *Enum) Values(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.options {
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(2), expose, s.Block(",")),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// ValueToStringCases generates code snippet cases from enum value to string
func (e *Enum) ValueToStringCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.options {
		name := v.name
		if len(v.enum) > 0 && e.storage == "enum" {
			name = v.enum
		}
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", name)),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// StringToValueCases generates code snippet cases from string to const value
func (e *Enum) StringToValueCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)

	for _, o := range e.options {
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.BlockF("case %s:", strings.Join(o.quotes(), ", "))),
			s.Compose(s.Indent(2), s.Block("return "), s.Block(o.value.Name()), s.Block(", nil")),
		)
	}

	return s.Snippets(s.NewLine(1), ss...)
}

// ValueToTextCases generates code snippet cases from enum value to text
func (e *Enum) ValueToTextCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.options {
		text := v.text
		if len(text) == 0 || text == v.value.Name() {
			text = v.name
		}
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", text)),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (e *Enum) ExtendKeys() []string {
	extends := slices.Collect(maps.Keys(e.extends))
	sort.Strings(extends)
	return extends
}

func (e *Enum) ExtendAttributes(ctx context.Context, name string) s.Snippet {
	ss := make([]s.Snippet, 0)

	ss = append(
		ss,
		s.Comments(fmt.Sprintf("%s describes %s attribute", name, name)),
		s.Compose(s.Block("func (v "), s.IdentTT(ctx, e.typ), s.BlockF(") %s() string {", name)),
		s.Compose(s.Indent(1), s.Block("switch v {")),
	)

	for _, v := range e.options {
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), s.BlockF("return %q", v.extends[name])),
		)
	}

	ss = append(
		ss,
		s.Compose(s.Indent(1), s.Block("default:")),
		s.Compose(s.Indent(2), s.BlockF("return %q", "")),
		s.Compose(s.Indent(1), s.Block("}")),
		s.Compose(s.Block("}\n")),
	)

	return s.Snippets(s.NewLine(1), ss...)
}

func (e *Enum) MappingKeys() []string {
	mappings := slices.Collect(maps.Keys(e.mapping))
	sort.Strings(mappings)
	return mappings
}

func (e *Enum) MappingAttributes(ctx context.Context, name string) s.Snippet {
	f := func(kind, value string) s.Snippet {
		if len(value) == 0 {
			return s.BlockF("return *new(%s)", kind)
		}
		return s.BlockF("return %s(%s)", kind, value)
	}

	ss := make([]s.Snippet, 0)
	kind := e.mapping[name]

	ss = append(
		ss,
		s.Comments(fmt.Sprintf("%s describes %s attribute", name, kind)),
		s.Compose(s.Block("func (v "), s.IdentTT(ctx, e.typ), s.BlockF(") %s() %s {", name, kind)),
		s.Compose(s.Indent(1), s.Block("switch v {")),
	)

	for _, v := range e.options {
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), f(kind, v.mapping[v.name])),
		)
	}

	ss = append(
		ss,
		s.Compose(s.Indent(1), s.Block("default:")),
		s.Compose(s.Indent(2), f(kind, "")),
		s.Compose(s.Indent(1), s.Block("}")),
		s.Compose(s.Block("}\n")),
	)

	return s.Snippets(s.NewLine(1), ss...)
}

func NewEnums(g genx.Context) *Enums {
	es := &Enums{
		e: make(map[types.Type]*Enum),
		p: g.Package(),
	}

	// Elements has been ordered by node(position)
	for elem := range es.p.Constants().Elements() {
		typ := elem.Type()
		if _, ok := typ.(*types.Named); !ok {
			continue
		}
		if elem.Value().Kind() != constant.Int {
			continue
		}

		if _, ok := es.e[typ]; !ok {
			v := &Enum{
				typ:     typ,
				key:     elem.TypeName(),
				options: make([]*option, 0),
				mapping: make(map[string]string),
				extends: make(map[string]bool),
			}
			es.e[typ] = v

			x := es.p.TypeNames().ElementByName(elem.TypeName())
			doc := docx.Parse(x.Doc())

			if annotations, ok := doc.AnnotationsByName(identifier); ok {
				for _, anno := range annotations {
					method, key, found := strings.Cut(anno.Key(), ".")
					if found {
						switch strings.ToLower(method) {
						case "map":
							v.mapping[key] = anno.Value()
						}
					} else {
						switch method {
						case "storage":
							v.storage = strings.ToLower(anno.Value())
						}
					}
				}
			}
		}
		es.e[typ].add(elem)
	}
	return es
}

type Enums struct {
	p pkgx.Package
	e map[types.Type]*Enum
}

func (es *Enums) Resolve(t types.Type) (*Enum, bool) {
	if _, ok := t.(*types.Named); !ok {
		return nil, false
	}
	e, ok := es.e[t]
	return e, ok
}
