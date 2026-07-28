package snippet

import (
	"context"
	"go/ast"
	"go/types"
	"iter"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
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

	p := dumper.Load(ctx, path)
	target := p.Types.Scope().Lookup(name)
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
			if tps, ok := x.Type().(interface {
				TypeParams() *types.TypeParamList
			}); ok {
				params = tps.TypeParams()
			}
		}
		if targc := params.Len(); targc != 0 {
			must.BeTrueF(
				targc == len(targs),
				"expected %d type parameter(s) for %s but got %d",
				targc, x.Name(), len(targs),
			)
			for i, targ := range targs {
				must.BeTrueF(
					!IsNil(targ),
					"got invalid type arg snippet at %d", i,
				)
				// must.BeTrueF(
				// 	ok,
				// 	"*types.Func type arguments MUST be a ident, but got %d:%T",
				// 	i, targ,
				// )
				r.targs = append(r.targs, targ)
			}
			// TODO should here need check the instantiation must can be succeeded by targs...
		}
		r.path = x.Pkg().Path()
		r.name = x.Name()
	case *types.Var, *types.Const:
		r.path = x.Pkg().Path()
		r.name = x.Name()
	}

	dumper.From(ctx).Track(ctx, path)
	return r
}

// ExposeUnsafe like Expose but skip validation of package and type arguments
func ExposeUnsafe(ctx context.Context, path, name string) Snippet {
	must.BeTrueF(path != "" && name != "", "package path and exposed name is required")
	must.BeTrueF(ast.IsExported(name), "exposed name must is exported")

	dumper.From(ctx).Track(ctx, path)
	return &exposer{
		path: path,
		name: name,
	}
}

// ExposeObject creates an exposer Snippet from a types.Object, validating the package and type arguments.
func ExposeObject(ctx context.Context, o types.Object, targs ...Snippet) Snippet {
	return Expose(ctx, o.Pkg().Path(), o.Name(), targs...)
}

// ExposeObjectUnsafe creates an exposer Snippet from a types.Object without validating the package or type arguments.
func ExposeObjectUnsafe(ctx context.Context, o types.Object) Snippet {
	return ExposeUnsafe(ctx, o.Pkg().Path(), o.Name())
}

/*
func ExposeByID(ctx context.Context, id string) Snippet {
	t := typx.LitTypeByID(id)

	targs := make([]Snippet, 0)
	for _, targ := range t.TypeArgs() {
		targs = append(targs, ExposeByID(ctx, targ.String()))
	}

	return &exposer{
		path:  t.PkgPath(),
		name:  t.Typename(),
		targs: targs,
	}
}
*/

type exposer struct {
	path  string
	name  string
	targs []Snippet
}

func (r *exposer) IsNil() bool {
	return false
}

func (r *exposer) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		b := &strings.Builder{}

		name := ""
		if r.path != "" {
			if namer, ok := typx.CtxPkgNamer.From(ctx); ok {
				name = namer.PackageName(r.path)
			}
		}
		if name != "" {
			b.WriteString(name + ".")
		}
		b.WriteString(r.name)
		if len(r.targs) > 0 {
			b.WriteString("[")
			for i, arg := range r.targs {
				if i > 0 {
					b.WriteString(", ")
				}
				for s := range arg.Fragments(ctx) {
					b.WriteString(s)
				}
			}
			b.WriteString("]")
		}

		yield(b.String())
	}
}
