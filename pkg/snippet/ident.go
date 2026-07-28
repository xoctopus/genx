package snippet

import (
	"context"
	"go/types"
	"iter"
	"reflect"

	"github.com/xoctopus/typx/pkg/typx"

	"github.com/xoctopus/genx/internal/dumper"
)

// IdentFor returns a Snippet representing the identifier for the generic type T.
func IdentFor[T any](ctx context.Context) Snippet {
	return _ident(ctx, typx.LitType(reflect.TypeFor[T]()))
}

// IdentOf returns a Snippet representing the identifier for the type of the given value v.
func IdentOf(ctx context.Context, v any) Snippet {
	return _ident(ctx, typx.LitType(reflect.TypeOf(v)))
}

// Ident returns a Snippet representing the identifier for the given typx.Type.
func Ident(ctx context.Context, t typx.Type) Snippet {
	return _ident(ctx, typx.LitType(t.Unwrap()))
}

// IdentRT returns a Snippet representing the identifier for the given reflect.Type.
func IdentRT(ctx context.Context, t reflect.Type) Snippet {
	return _ident(ctx, typx.LitType(t))
}

// IdentTT returns a Snippet representing the identifier for the given types.Type.
func IdentTT(ctx context.Context, t types.Type) Snippet {
	return _ident(ctx, typx.LitType(t))
}

func _ident(ctx context.Context, t *typx.Literal) Snippet {
	tracker := dumper.From(ctx)

	var walk func(lit *typx.Literal)
	walk = func(lit *typx.Literal) {
		if lit == nil {
			return
		}
		tracker.Track(ctx, lit.PkgPath())
		for _, targ := range lit.TypeArgs() {
			walk(targ)
		}
		walk(lit.Key())
		walk(lit.Elem())
		for _, in := range lit.Ins() {
			walk(in)
		}
		for _, out := range lit.Outs() {
			walk(out)
		}
		if lit.Typename() == "" {
			for _, out := range lit.Fields() {
				walk(out)
			}
			for _, method := range lit.Methods() {
				walk(method)
			}
		}
	}
	walk(t)
	return &ident{t: t}
}

type ident struct {
	t *typx.Literal
}

func (v *ident) IsNil() bool {
	return v == nil || v.t == nil
}

func (v *ident) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		yield(v.t.Dump(ctx))
	}
}
