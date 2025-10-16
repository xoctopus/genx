package genx

import (
	"context"
	"go/token"
	"go/types"

	"github.com/xoctopus/pkgx"

	"github.com/xoctopus/genx/snippet"
)

type GeneratorNewer interface {
	Newer(Context) Generator
}

type Generator interface {
	Identifier() string
	Generate(Context, types.Type) error
}

type GeneratorPoster interface {
	Post() []byte
}

type Context interface {
	IsZero() bool

	Package() pkgx.Package
	PackageByPath(string) pkgx.Package
	PackageByPos(token.Pos) pkgx.Package

	Render(snippet snippet.Snippet)
}

type Executor interface {
	Execute(context.Context, ...Generator) error
}

type Args struct {
	Entrypoint     []string
	FilenameSuffix string
}

func NewContext(args *Args) Executor {
	return &genc{
		args: args,
		u:    pkgx.NewPackages(args.Entrypoint...),
	}
}

type genc struct {
	args *Args

	u *pkgx.Packages
	p *pkgx.Package
}

func (x *genc) IsZero() bool {
	return false
}

func (x *genc) Execute(context.Context, ...Generator) error {
	return nil
}
