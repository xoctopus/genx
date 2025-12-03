package codex

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"

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
