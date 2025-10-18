package errorx

import (
	"context"
	"fmt"
	"go/constant"
	"go/types"
	"strings"

	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

type Error struct {
	name      string
	typ       types.Type
	def       string
	undefined *pkgx.Constant
	list      []*pkgx.Constant
}

func (e *Error) IsValid() bool {
	return e.undefined != nil && len(e.list) > 0
}

func (e *Error) add(c *pkgx.Constant) {
	name := c.Name()

	if name[0] == '_' {
		return
	}

	prefix := stringsx.UpperSnakeCase(e.name)
	if name == prefix+"_UNDEFINED" {
		e.undefined = c
		return
	}

	parts := strings.SplitN(name, "__", 2)
	if len(parts) == 2 && parts[0] == prefix {
		e.list = append(e.list, c)
	}
}

// MessageKeyValues generates code snippet error message key value pairs
func (e *Error) MessageKeyValues(ctx context.Context) s.Snippet {
	ss := []s.Snippet{
		s.Compose(
			s.Indent(1),
			s.ExposeObject(ctx, e.undefined.Exposer()),
			s.BlockF(": %q,", fmt.Sprintf("[%s:%s] undefined", e.def, e.undefined.Value())),
		),
	}
	for _, v := range e.list {
		msg := strings.Join(v.Doc().Desc(), " ")
		if len(msg) == 0 {
			msg = strings.TrimPrefix(v.Name(), stringsx.UpperSnakeCase(v.TypeName())+"__")
		}
		ss = append(
			ss,
			s.Compose(
				s.Indent(1),
				s.ExposeObject(ctx, v.Exposer()),
				s.Block(": "),
				s.BlockRaw(fmt.Sprintf("[%s:%s] %s", e.def, v.Value(), msg)),
				s.Block(","),
			),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func NewErrors(g genx.Context) *Errors {
	es := &Errors{
		p: g.Package(),
		e: make(map[types.Type]*Error),
	}

	for elem := range es.p.Constants().Elements() {
		typ := elem.Type()
		if _, ok := typ.(*types.Named); !ok {
			continue
		}
		if elem.Value().Kind() != constant.Int {
			continue
		}

		if _, ok := es.e[typ]; !ok {
			te := es.p.TypeNames().ElementByName(elem.TypeName())
			must.BeTrue(te != nil)
			def := ""
			for _, desc := range te.Doc().Desc() {
				if strings.HasPrefix(desc, "@def ") {
					def = strings.TrimPrefix(desc, "@def ")
					break
				}
			}
			if def == "" {
				def = es.p.Unwrap().Name() + "." + elem.TypeName()
			}
			es.e[typ] = &Error{
				typ:  typ,
				def:  def,
				name: elem.TypeName(),
				list: make([]*pkgx.Constant, 0),
			}
		}
		es.e[typ].add(elem)
	}
	return es
}

type Errors struct {
	p pkgx.Package
	e map[types.Type]*Error
}

func (es *Errors) Resolve(t types.Type) (*Error, bool) {
	if _, ok := t.(*types.Named); !ok {
		return nil, false
	}
	e, ok := es.e[t]
	return e, ok
}
