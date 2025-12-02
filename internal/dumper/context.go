package dumper

import (
	"context"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/contextx"
)

var ctxTracker = contextx.NewT[ImportTracker]()

func TrackerFrom(ctx context.Context) ImportTracker {
	return ctxTracker.MustFrom(ctx)
}

func WithTracker(ctx context.Context, path, module string) context.Context {
	if _, ok := ctxTracker.From(ctx); ok {
		return ctx
	}
	i := NewImportTracker(path, module)
	return ctxTracker.With(typx.CtxPkgNamer.With(ctx, i), i)
}

func TrackerCarrier(path, module string) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return WithTracker(ctx, path, module)
	}
}
