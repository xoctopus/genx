package dumper

import (
	"context"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/contextx"
)

var ctxTracker = contextx.NewT[ImportTracker]()

func From(ctx context.Context) ImportTracker {
	return ctxTracker.MustFrom(ctx)
}

func With(ctx context.Context, t ImportTracker) context.Context {
	ctx = typx.CtxPkgNamer.With(ctx, t)
	return ctxTracker.With(ctx, t)
}

func WithEntry(ctx context.Context, entry string) context.Context {
	return With(ctx, NewImportTracker(entry))
}

func Carrier(t ImportTracker) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, t)
	}
}

func CarrierEntry(entry string) contextx.Carrier {
	return Carrier(NewImportTracker(entry))
}
