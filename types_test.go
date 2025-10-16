package genx_test

import (
	"context"
	"go/types"
	"strconv"

	"github.com/xoctopus/genx"
	s "github.com/xoctopus/genx/snippet"
)

type TestGenerator struct {
}

func (g *TestGenerator) Identifier() string {
	return "test_genx"
}

func (g *TestGenerator) Generate(c genx.Context, t types.Type) error {
	tx, ok := t.(*types.Named)
	if !ok {
		return nil
	}

	if c.Package().Unwrap().Path() != "" {

	}

	var obj *types.TypeName
	for exp := range c.Package().TypeNames().Exposers() {
		if types.Identical(exp.Type(), tx) {
			if exp.Name() == "Gender" {
				obj = exp
				break
			}
		}
	}
	if obj == nil {
		return nil
	}

	ctx := c.Context()

	c.Render(s.Snippets(
		s.NewLine(1),
		s.Compose(
			s.Block("var _ = "),
			s.ExposeObject(ctx, obj),
			s.Block("(1)"),
		),
		s.Compose(
			s.Block("var _ = "),
			s.Expose(ctx, "github.com/pkg/errors", "New"),
			s.BlockF("(%s)", strconv.Quote("some error")),
		),
		s.Compose(
			s.Block("var _ = "),
			s.Expose(ctx, "bytes", "NewBufferString"),
			s.BlockF("(%s)", strconv.Quote("some string for bytes.Buffer")),
		),
	))
	return nil
}

func ExampleGenerator() {
	c := genx.NewContext(&genx.Args{
		Entrypoint: []string{"github.com/xoctopus/genx/testdata"},
	})
	_ = c.Execute(context.Background(), &TestGenerator{})

	// Output:

}
