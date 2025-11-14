package dumper

import (
	"context"

	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/contextx"
)

var ctx = contextx.NewT[ImportTracker]()

func TrackerFromContext(child context.Context) ImportTracker {
	return ctx.MustFrom(child)
}

func WithTrackerContext(parent context.Context, path, module string) context.Context {
	if _, ok := ctx.From(parent); ok {
		return parent
	}

	i := NewImportTracker(path, module)
	return ctx.With(pkgx.WithNamer(parent, i), i)
}
