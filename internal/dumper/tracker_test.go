package dumper_test

import (
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

	tracker.Track("github.com/xoctopus/genx/testdata")
	tracker.Track("github.com/pkg/errors")
	tracker.Track("errors")
	tracker.Track("bytes")
	tracker.Track("strings")
	tracker.Track("io")
	tracker.Track("errors")
	tracker.Track("fmt")
	tracker.Track("context")
	tracker.Track("")                                  // track empty
	tracker.Track("github.com/xoctopus/genx/testdata") // track tracked
	tracker.Track("github.com/xoctopus/pkgx")
	tracker.Track("github.com/xoctopus/typex")

	tracker.Init()

	t.Run("TrackAfterInitialized", func(t *testing.T) {
		ExpectPanic(t, func() { tracker.Track("any") }, ErrorContains("cannot track"))
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
