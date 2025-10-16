package dumper

import (
	"context"

	"github.com/xoctopus/typex/namer"
	"github.com/xoctopus/x/contextx"
)

var ctx = contextx.NewT[ImportTracker]()

func TrackerFromContext(child context.Context) ImportTracker {
	return ctx.MustFrom(child)
}

func WithTrackerContext(parent context.Context, entry string) context.Context {
	if _, ok := ctx.From(parent); ok {
		return parent
	}

	i := NewImportTracker(entry)
	return ctx.With(namer.WithContext(parent, i), i)
}

func Track(child context.Context, path string) {
	ctx.MustFrom(child).Track(path)
}
