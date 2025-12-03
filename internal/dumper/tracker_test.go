package dumper_test

import (
	"context"
	"strconv"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestNewImportTracker(t *testing.T) {
	tracker := dumper.NewImportTracker("github.com/xoctopus/genx/testdata")

	Expect(t, tracker.Entry(), Equal("github.com/xoctopus/genx/testdata"))
	Expect(t, tracker.Module(), Equal("github.com/xoctopus/genx"))

	t.Run("Tracking", func(t *testing.T) {
		t.Run("FetchBeforeInitialized", func(t *testing.T) {
			ExpectPanic(t, func() { tracker.PackageName("any") }, ErrorContains("cannot fetch"))
			ExpectPanic(t, func() { tracker.Range(nil) }, ErrorContains("cannot range"))
		})

		ctx := context.Background()
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata")
		tracker.Track(ctx, "errors")
		tracker.Track(ctx, "bytes")
		tracker.Track(ctx, "strings")
		tracker.Track(ctx, "io")
		tracker.Track(ctx, "errors")
		tracker.Track(ctx, "fmt")
		tracker.Track(ctx, "context")
		tracker.Track(ctx, "")                                  // track empty
		tracker.TrackCustom(ctx, "", "_")                       // track empty
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata") // track tracked
		tracker.Track(ctx, "github.com/xoctopus/pkgx/pkg/pkgx")
		tracker.Track(ctx, "github.com/xoctopus/typx/pkg/typx")
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata/errors") // conflict with std.errors
		// use custom package alias. duplicately import will not take effect
		tracker.TrackCustom(ctx, "github.com/xoctopus/typx/pkg/typx", "_")
		// use custom package alias.
		tracker.TrackCustom(ctx, "github.com/xoctopus/typx/internal/dumper", "_")
		// conflict package names
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata/bricks/x/p")
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata/bricks/y/p")
		tracker.Track(ctx, "github.com/xoctopus/genx/testdata/bricks/z/sub/p")

		tracker.Init()
		t.Run("TrackAfterInitialized", func(t *testing.T) {
			ExpectPanic(t, func() { tracker.Track(ctx, "any") }, ErrorContains("cannot track"))
		})

		t.Run("IsEntry", func(t *testing.T) {
			Expect(t, tracker.PackageName("github.com/xoctopus/genx/testdata"), HaveLen[string](0))
		})
		t.Run("Unimported", func(t *testing.T) {
			ExpectPanic(t, func() { tracker.PackageName("unimported") }, ErrorContains("not be tracked"))
		})
		t.Run("Namer", func(t *testing.T) {
			Expect(t, tracker.PackageName("github.com/xoctopus/genx/testdata/errors"), Equal("testdata_errors"))
		})

		imports := make([]string, 0)
		for i := range tracker.Range {
			imports = append(imports, i.Code())
			if i.Code() == strconv.Quote(i.Path()) {
				Expect(t, i.Name(), Equal(i.Alias()))
			} else {
				Expect(t, i.Name(), NotEqual(i.Alias()))
			}
		}
		Expect(t, imports, Equal([]string{
			`"bytes"`,
			`"context"`,
			`"errors"`,
			`"fmt"`,
			`"io"`,
			`"strings"`,
			`x_p "github.com/xoctopus/genx/testdata/bricks/x/p"`,
			`y_p "github.com/xoctopus/genx/testdata/bricks/y/p"`,
			`sub_p "github.com/xoctopus/genx/testdata/bricks/z/sub/p"`,
			`testdata_errors "github.com/xoctopus/genx/testdata/errors"`,
			`"github.com/xoctopus/pkgx/pkg/pkgx"`,
			`_ "github.com/xoctopus/typx/internal/dumper"`,
			`"github.com/xoctopus/typx/pkg/typx"`,
		}))
	})
}
