package dumper_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestWithTrackerContext(t *testing.T) {
	ctx := context.Background()

	ctx = dumper.WithTracker(ctx, "pkg/path", "pkg/module")
	ctx2 := dumper.TrackerCarrier("any", "any")(ctx)

	Expect(t, ctx, Equal(ctx2))
	Expect(t, dumper.TrackerFrom(ctx), Equal(dumper.TrackerFrom(ctx2)))
}
