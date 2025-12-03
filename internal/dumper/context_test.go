package dumper_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestWithTrackerContext(t *testing.T) {
	ctx := context.Background()

	ctx = dumper.WithTracker(ctx, "github.com/xoctopus/genx/testdata")
	Expect(t, dumper.TrackerFrom(ctx).Entry(), Equal("github.com/xoctopus/genx/testdata"))

	ctx = dumper.TrackerCarrier("github.com/xoctopus/genx/testdata")(ctx)
	Expect(t, dumper.TrackerFrom(ctx).Entry(), Equal("github.com/xoctopus/genx/testdata"))
}
