package enumx

import (
	"context"
	"fmt"
	"go/constant"
	"go/types"
	"sort"
	"strings"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/pkg/docx"
	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

type option struct {
	name  string
	text  string
	attrs map[string]string
	value *pkgx.Constant
}

type Enum struct {
	typ     types.Type
	key     string
	unknown *pkgx.Constant
	values  []*option
	attrs   map[string]struct{}
	options map[string]string
	storage string
}

func (e *Enum) IsValid() bool {
	return e.unknown != nil || len(e.values) > 0
}

// add adds option
func (e *Enum) add(c *pkgx.Constant) {
	if name := c.Name(); name != "_" {
		prefix := stringsx.UpperSnakeCase(e.key)
		if name == prefix+"_UNKNOWN" {
			e.unknown = c
			return
		}

		before, after, found := strings.Cut(name, "__")
		if !found || before != prefix {
			return
		}

		o := &option{
			value: c,
			name:  after,
			text:  "",
			attrs: map[string]string{},
		}

		doc := docx.Parse(c.Doc())
		for key, annotations := range doc.Annotations() {
			for _, anno := range annotations {
				switch key {
				case "attr":
					if k, v, found := strings.Cut(anno.Text(), "="); found {
						o.attrs[k] = v
						e.attrs[k] = struct{}{}
					}
				default:
				}
			}
		}
		o.text = doc.Title(c.Name())
		e.values = append(e.values, o)
	}
}

// Values generates code snippet of const value list
func (e *Enum) Values(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
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
	for _, v := range e.values {
		name := strings.TrimPrefix(
			v.value.Name(),
			stringsx.UpperSnakeCase(v.value.TypeName())+"__",
		)
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
	for _, v := range e.values {
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.BlockF("case %q:", v.name)),
			s.Compose(s.Indent(2), s.Block("return "), expose, s.Block(", nil")),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

// ValueToTextCases generates code snippet cases from enum value to text
func (e *Enum) ValueToTextCases(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, v := range e.values {
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

func (e *Enum) Attr(ctx context.Context, attr, option string) s.Snippet {
	f := func(v string) s.Snippet {
		switch option {
		case "string", "text", "":
			return s.BlockF("return %q", v)
		default:
			if len(v) == 0 {
				return s.BlockF("return *new(%s)", option)
			}
			return s.BlockF("return %s(%s)", option, v)
		}
	}

	if option == "" {
		option = "string"
	}

	ss := make([]s.Snippet, 0)

	name := stringsx.UpperCamelCase(attr)

	ss = append(
		ss,
		s.Comments(fmt.Sprintf("%s describes %s attribute", name, name)),
		s.Compose(s.Block("func (v "), s.IdentTT(ctx, e.typ), s.BlockF(") %s() %s {", name, option)),
		s.Compose(s.Indent(1), s.Block("switch v {")),
	)

	for _, v := range e.values {
		expose := s.ExposeObjectUnsafe(ctx, v.value.Exposer())
		ss = append(
			ss,
			s.Compose(s.Indent(1), s.Block("case "), expose, s.Block(":")),
			s.Compose(s.Indent(2), f(v.attrs[attr])),
		)
	}

	ss = append(
		ss,
		s.Compose(s.Indent(1), s.Block("default:")),
		s.Compose(s.Indent(2), f("")),
		s.Compose(s.Indent(1), s.Block("}")),
		s.Compose(s.Block("}\n")),
	)

	return s.Snippets(s.NewLine(1), ss...)
}

func (e *Enum) Attrs() []string {
	attrs := make([]string, 0, len(e.attrs))
	for attr := range e.attrs {
		attrs = append(attrs, attr)
	}
	sort.Strings(attrs)
	return attrs
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
				values:  make([]*option, 0),
				attrs:   make(map[string]struct{}),
				options: make(map[string]string),
			}
			es.e[typ] = v

			x := es.p.TypeNames().ElementByName(elem.TypeName())
			d := docx.Parse(x.Doc())

			if annotations, ok := d.AnnotationsByName("def"); ok {
				for _, anno := range annotations {
					if key := anno.Key(); len(key) > 0 {
						if key == "storage" && len(v.storage) == 0 {
							v.storage = strings.ToLower(anno.Value())
						} else {
							if option, ok := strings.CutPrefix(key, "attr."); ok {
								v.options[option] = anno.Value()
							}
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
