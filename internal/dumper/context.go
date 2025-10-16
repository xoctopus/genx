package dumper

import (
	"context"

	"github.com/xoctopus/typex/namer"
	"github.com/xoctopus/x/contextx"
)

var (
	trackerCtx = contextx.NewT[ImportTracker]()
	selfCtx    = contextx.NewT[string]()
)

func TrackerFromContext(ctx context.Context) ImportTracker {
	return trackerCtx.MustFrom(ctx)
}

func WithTrackerContext(ctx context.Context) context.Context {
	if _, ok := trackerCtx.From(ctx); ok {
		return ctx
	}

	i := NewImportTracker(SelfFromContext(ctx))
	return trackerCtx.With(namer.WithContext(ctx, i), i)
}

func SelfFromContext(ctx context.Context) string {
	v, _ := selfCtx.From(ctx)
	return v
}

func WithSelfContext(ctx context.Context, path string) context.Context {
	return selfCtx.With(ctx, path)
}
