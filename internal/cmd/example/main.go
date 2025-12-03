package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/genx/devpkg/codex"
	_ "github.com/xoctopus/genx/devpkg/enumx"
	"github.com/xoctopus/genx/pkg/genx"
)

func main() {
	entry := filepath.Join(must.NoErrorV(os.Getwd()), "devpkg", "testdata")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{
			filepath.Join(entry, "enumx"),
			filepath.Join(entry, "codex"),
		},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}
}
