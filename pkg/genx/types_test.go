package genx_test

import (
	"context"
	"go/types"
	"strconv"

	"github.com/xoctopus/pkgx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

type TestGenerator struct {
}

func (g *TestGenerator) Identifier() string {
	return "test_genx"
}

func (g *TestGenerator) Generate(c genx.Context, t types.Type) error {
	must.BeTrue(!c.IsZero())
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

type TestGeneratorHasSyntaxError struct {
	p pkgx.Package
}

func (v *TestGeneratorHasSyntaxError) Identifier() string {
	return "test_genx_e"
}

func (v *TestGeneratorHasSyntaxError) New(c genx.Context) genx.Generator {
	return &TestGeneratorHasSyntaxError{p: c.PackageByPath("errors")}
}

func (v *TestGeneratorHasSyntaxError) Generate(c genx.Context, t types.Type) error {
	must.BeTrue(!c.IsZero())
	c.Render(s.Snippets(
		s.NewLine(1),
		s.Block("func X() int {"),
		s.Block("\treturn 1"),
	))
	return nil
}

type TestGeneratorHasTypeGenerated struct{}

func (v *TestGeneratorHasTypeGenerated) Identifier() string {
	return "test_genx_t"
}

func (v *TestGeneratorHasTypeGenerated) Generate(c genx.Context, _ types.Type) error {
	must.BeTrue(!c.IsZero())
	ctx := c.Context()
	c.Render(s.Snippets(
		s.NewLine(1),
		s.Block("type ("),
		s.Compose(s.Indent(1), s.Block("NetAddr = "), s.Expose(ctx, "net", "Addr")),
		s.Compose(s.Indent(1), s.Block("Buffer "), s.Expose(ctx, "bytes", "Buffer")),
		s.Block(")"),
	))
	return nil
}

func ExampleGenerator() {
	c := genx.NewContext(&genx.Args{
		Entrypoint: []string{"github.com/xoctopus/genx/testdata"},
	})
	_ = c.Execute(context.Background(), &TestGenerator{}, &TestGeneratorHasTypeGenerated{})

	c = genx.NewContext(&genx.Args{
		Entrypoint: []string{"github.com/xoctopus/genx/testdata/ignored"},
	})
	_ = c.Execute(context.Background(), &TestGenerator{})

	c = genx.NewContext(&genx.Args{
		Entrypoint: []string{"github.com/xoctopus/genx/testdata/ignored"},
	})
	_ = c.Execute(context.Background(), &TestGeneratorHasSyntaxError{})

	c = genx.NewContext(&genx.Args{
		Entrypoint: []string{"github.com/xoctopus/genx/testdata"},
	})
	_ = c.Execute(context.Background(), &TestGeneratorHasSyntaxError{})

	//Output:
	//    2: package testdata
	//    3:
	//    4:
	//    5: func X() int {
	//    6: 	return 1
	//                ↑
	//expected ';', found 'EOF'
}
