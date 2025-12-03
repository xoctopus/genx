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

func WithTracker(ctx context.Context, entry string) context.Context {
	i := NewImportTracker(entry)
	return ctxTracker.With(typx.CtxPkgNamer.With(ctx, i), i)
}

func TrackerCarrier(entry string) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return WithTracker(ctx, entry)
	}
}
