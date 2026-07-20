package contextx_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/genx/devpkg/contextx"
	"github.com/xoctopus/genx/pkg/genx"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())
	entry := filepath.Join(cwd, "..", "testdata", "contextx")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})
	must.NoError(ctx.Execute(context.Background(), genx.Get()...))

	// Output:
}
