package snippet

import (
	"context"
	"go/types"
	"iter"
	"reflect"

	"github.com/xoctopus/typx/pkg/typx"

	"github.com/xoctopus/genx/internal/dumper"
)

func IdentFor[T any](ctx context.Context) Snippet {
	return Ident(ctx, typx.NewRType(reflect.TypeFor[T]()))
}

func IdentOf[T any](ctx context.Context, v T) Snippet {
	return Ident(ctx, typx.NewRType(reflect.TypeOf(v)))
}

func Ident(ctx context.Context, t typx.Type) Snippet {
	dumper.TrackerFrom(ctx).Track(ctx, t.PkgPath())
	return &ident{t: t}
}

func IdentRT(ctx context.Context, t reflect.Type) Snippet {
	return Ident(ctx, typx.NewRType(t))
}

func IdentTT(ctx context.Context, t types.Type) Snippet {
	return Ident(ctx, typx.NewTType(t))
}

type ident struct {
	t typx.Type
}

func (v *ident) IsNil() bool {
	return false
}

func (v *ident) Fragments(ctx context.Context) iter.Seq[string] {
	return func(yield func(string) bool) {
		yield(typx.TypeLit(ctx, v.t.Unwrap()))
	}
}
