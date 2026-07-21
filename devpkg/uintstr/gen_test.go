package uintstr_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/genx/devpkg/uintstr"
	"github.com/xoctopus/genx/pkg/genx"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())

	entry := filepath.Join(cwd, "..", "testdata", "uintstr")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}

	// Output:
}
