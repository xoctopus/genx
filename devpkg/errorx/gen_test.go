package errorx

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/genx/pkg/genx"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())
	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{filepath.Join(cwd, "testdata")},
	})
	must.NoError(ctx.Execute(context.Background(), genx.Get()...))

	// Output:
}
