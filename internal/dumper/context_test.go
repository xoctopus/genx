package dumper_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestWithTrackerContext(t *testing.T) {
	ctx := context.Background()

	ctx = dumper.WithTrackerContext(ctx, "pkg/path", "pkg/module")
	ctx2 := dumper.WithTrackerContext(ctx, "any", "any")

	Expect(t, ctx, Equal(ctx2))
	Expect(t, dumper.TrackerFromContext(ctx), Equal(dumper.TrackerFromContext(ctx2)))
}
