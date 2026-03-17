package docx_test

import (
	"context"
	"os"
	"path/filepath"

	_ "github.com/xoctopus/genx/devpkg/docx"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/genx/pkg/genx"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())

	entry := filepath.Join(cwd, "..", "testdata", "docx")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	must.NoError(ctx.Execute(context.Background(), genx.Get()...))
	// Output:
}
