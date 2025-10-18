package errorx

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/stringsx"

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

func TestX(t *testing.T) {
	t.Log(stringsx.UpperSnakeCase("Code3"))
}
