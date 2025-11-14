package dumper_test

import (
	"context"
	"strconv"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestNewImportTracker(t *testing.T) {
	path := "github.com/xoctopus/genx/testdata"
	module := "github.com/xoctopus/genx/testdata"

	tracker := dumper.NewImportTracker(path, module)

	Expect(t, tracker.Entry(), Equal(path))
	Expect(t, tracker.Module(), Equal(path))

	t.Run("FetchBeforeInitialized", func(t *testing.T) {
		ExpectPanic(t, func() { tracker.Package("any") }, ErrorContains("cannot fetch"))
		ExpectPanic(t, func() { tracker.Range(nil) }, ErrorContains("cannot range"))
	})

	ctx := context.Background()

	tracker.Track(ctx, "github.com/xoctopus/genx/testdata")
	tracker.Track(ctx, "github.com/pkg/errors")
	tracker.Track(ctx, "errors")
	tracker.Track(ctx, "bytes")
	tracker.Track(ctx, "strings")
	tracker.Track(ctx, "io")
	tracker.Track(ctx, "errors")
	tracker.Track(ctx, "fmt")
	tracker.Track(ctx, "context")
	tracker.Track(ctx, "")                                  // track empty
	tracker.Track(ctx, "github.com/xoctopus/genx/testdata") // track tracked
	tracker.Track(ctx, "github.com/xoctopus/pkgx")
	tracker.Track(ctx, "github.com/xoctopus/typex")

	tracker.Init()

	t.Run("TrackAfterInitialized", func(t *testing.T) {
		ExpectPanic(t, func() { tracker.Track(ctx, "any") }, ErrorContains("cannot track"))
	})

	Expect(t, tracker.Package("github.com/xoctopus/genx/testdata"), HaveLen[string](0))
	ExpectPanic(t, func() { tracker.Package("unimported") }, ErrorContains("not be tracked"))
	Expect(t, tracker.Package("github.com/pkg/errors"), Equal("pkg_errors"))

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
		`pkg_errors "github.com/pkg/errors"`,
		`"github.com/xoctopus/pkgx"`,
		`"github.com/xoctopus/typex"`,
	}))
}
