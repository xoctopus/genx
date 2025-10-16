package snippet

import (
	"context"
	"go/ast"
	"go/types"
	"iter"

	"github.com/xoctopus/typex/namer"
	"github.com/xoctopus/typex/pkgutil"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/genx/internal/dumper"
)

// Expose create an exposer in some package, it may be a types.Type or a types.Object
// here the name MUST BE exported
// case 1: a named/alias type. this case should handled by Ident
// eg: path/to/package.NamedType[TypeArguments...]
// case 2: an exported object. the package level MUST be a *types.Func, *types.Const or *types.Var
// eg: errors.New the underlying is a function: func() error, but we need the object
func Expose(ctx context.Context, path string, name string, targs ...Snippet) Snippet {
	must.BeTrueF(
		path != "" && name != "",
		"package path and exposed name is required",
	)
	must.BeTrueF(
		ast.IsExported(name),
		"exposed name must is exported",
	)

	p := pkgutil.New(path)
	target := p.Scope().Lookup(name)
	must.BeTrueF(
		target != nil,
		"cannot lookup `%s` in package `%s`",
		name, path,
	)

	r := &exposer{}
	switch x := target.(type) {
	case *types.Func, *types.TypeName:
		var params *types.TypeParamList
		if _, ok := x.(*types.Func); ok {
			params = x.Type().(*types.Signature).TypeParams()
		} else {
			params = x.Type().(*types.Named).TypeParams()
		}
		if targc := params.Len(); targc != 0 {
			must.BeTrueF(
				targc == len(targs),
				"expected %d type parameter(s) for %s but got %d",
				targc, x.Name(), len(targs),
			)
			for i, targ := range targs {
				must.BeTrueF(
					targ != nil && !targ.IsNil(),
					"got invalid type arg snippet at %d", i,
				)
				ta, ok := targ.(*ident)
				must.BeTrueF(
					ok,
					"*types.Func type arguments MUST be a ident, but got %d:%T",
					i, targ,
				)
				r.targs = append(r.targs, ta)
			}
			// TODO should here need check the instantiation must can be succeeded by targs...
		}
		r.path = x.Pkg().Path()
		r.name = x.Name()
	case *types.Var, *types.Const:
		r.path = x.Pkg().Path()
		r.name = x.Name()
	}

	dumper.Track(ctx, path)
	return r
}

func ExposeObject(ctx context.Context, o types.Object, targs ...Snippet) Snippet {
	return Expose(ctx, o.Pkg().Path(), o.Name(), targs...)
}

type exposer struct {
	path  string
	name  string
	targs []*ident
}

func (r *exposer) IsNil() bool {
	return false
}

func (r *exposer) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		path := namer.MustFromContext(ctx).Package(r.path)
		if !yield(path) {
			return
		}
		if path != "" {
			if !yield(".") {
				return
			}
		}
		if !yield(r.name) {
			return
		}
		if len(r.targs) > 0 {
			if !yield("[") {
				return
			}
			for i, arg := range r.targs {
				if i > 0 {
					if !yield(", ") {
						return
					}
				}
				for s := range arg.Fragments(ctx) {
					if !yield(s) {
						return
					}
				}
			}
			if !yield("]") {
				return
			}
		}
	}
}
