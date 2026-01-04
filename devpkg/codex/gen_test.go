package codex

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/testx"

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
