package snippet

import (
	"context"
	"go/types"
	"iter"
	"reflect"

	"github.com/xoctopus/typex"

	"github.com/xoctopus/genx/internal/dumper"
)

func IdentFor[T any](ctx context.Context) Snippet {
	return Ident(ctx, typex.NewRType(ctx, reflect.TypeFor[T]()))
}

func IdentOf[T any](ctx context.Context, v T) Snippet {
	return Ident(ctx, typex.NewRType(ctx, reflect.TypeOf(v)))
}

func Ident(ctx context.Context, t typex.Type) Snippet {
	dumper.TrackerFrom(ctx).Track(ctx, t.PkgPath())
	return &ident{t: t}
}

func IdentRT(ctx context.Context, t reflect.Type) Snippet {
	return Ident(ctx, typex.NewRType(ctx, t))
}

func IdentTT(ctx context.Context, t types.Type) Snippet {
	return Ident(ctx, typex.NewTType(ctx, t))
}

type ident struct {
	t typex.Type
}

func (v *ident) IsNil() bool {
	return false
}

func (v *ident) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		yield(v.t.TypeLit(ctx))
	}
}
