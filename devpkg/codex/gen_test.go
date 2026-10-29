package codex

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/devpkg/testdata/codex"
	"github.com/xoctopus/genx/pkg/genx"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())
	entry := filepath.Join(cwd, "..", "testdata", "codex")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})
	must.NoError(ctx.Execute(context.Background(), genx.Get()...))

	// Output:
}

func TestG_Generate(t *testing.T) {
	gx := genx.Get("code")[0]
	testx.Expect(t, gx, testx.NotBeNil[genx.Generator]())
	testx.Expect(t, gx.Generate(nil, nil), testx.Succeed())
}

func TestGeneratedCode(t *testing.T) {
	var (
		_ codex.ErrorMessage = testdata.CODE__ERROR1
		_ codex.CanBeCode    = testdata.CODE__ERROR1
		_ codex.CanBeCode    = testdata.DOMAIN_CODE__PARSE
	)

	t.Run("CodeIdentifier", func(t *testing.T) {
		testx.Expect(t, testdata.CODE_UNDEFINED.CodeIdentifier(), testx.Equal("TESTDATA::CODE__UNDEFINED"))
		testx.Expect(t, testdata.CODE__ERROR1.CodeIdentifier(), testx.Equal("TESTDATA::CODE__ERROR1"))
		testx.Expect(t, testdata.DOMAIN_CODE__PARSE.CodeIdentifier(), testx.Equal("DOMAIN_NAME__PARSE"))
		testx.Expect(t, testdata.Code(100).CodeIdentifier(), testx.Equal(""))
	})

	t.Run("CodeHint", func(t *testing.T) {
		_, ok := any(testdata.CODE__ERROR1).(interface{ CodeHint() string })
		testx.Expect(t, ok, testx.BeFalse())

		testx.Expect(t, testdata.DOMAIN_CODE__PARSE.CodeHint(), testx.Equal("parse failed"))
		testx.Expect(t, testdata.DOMAIN_CODE__PARAMETER.CodeHint(), testx.Equal("please check the input parameters"))
		testx.Expect(t, testdata.DOMAIN_CODE__HANDLE.CodeHint(), testx.Equal(""))
		testx.Expect(t, testdata.DOMAIN_CODE__INTERNAL.CodeHint(), testx.Equal(""))
		testx.Expect(t, testdata.DOMAIN_CODE_UNDEFINED.CodeHint(), testx.Equal(""))
	})

	t.Run("CodeValues", func(t *testing.T) {
		testx.Expect(t, testdata.CODE_UNDEFINED.CodeValues(), testx.Equal([]any{
			testdata.CODE__ERROR1,
			testdata.CODE__ERROR2,
			testdata.CODE__ERROR3,
			testdata.CODE__ERROR4,
		}))
	})
}
