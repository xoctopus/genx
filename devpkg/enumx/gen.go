package enumx

import (
	"bytes"
	"database/sql/driver"
	_ "embed"
	"go/types"
	"log"
	"strings"
	"time"

	"github.com/xoctopus/x/enumx"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

var (
	//go:embed enumx.go.tpl
	template []byte
	//go:embed enumx.go_int.tpl
	templateInt []byte // template for int storage driver.Valuer/sql.Scanner
	//go:embed enumx.go_txt.tpl
	templateTxt []byte // template for text/varchar storage driver.Valuer/sql.Scanner
)

func init() {
	genx.Register(&g{})
}

type g struct {
	enums *Enums
}

func (x *g) Identifier() string {
	return "enum"
}

func (x *g) New(c genx.Context) genx.Generator {
	return &g{enums: NewEnums(c)}
}

func (x *g) Generate(c genx.Context, t types.Type) error {

	if e, ok := x.enums.Resolve(t); ok {
		if e.IsValid() {
			cost := Span()
			log.Printf("genx:%s %s\n", x.Identifier(), t.String())
			x.generate(c, e)
			log.Printf("==> cost: %fs", cost().Seconds())
			return nil
		}
	}
	return nil
}

func (x *g) generate(c genx.Context, e *Enum) {
	ctx := c.Context()

	ident := s.IdentTT(ctx, e.typ)
	pkgid := "github.com/xoctopus/x/enumx"

	args := []*s.TArg{
		// @def bytes.ToUpper
		s.ArgExposeUnsafe(ctx, "bytes", "ToUpper"),
		// @def fmt.Sprintf
		s.ArgExposeUnsafe(ctx, "fmt", "Sprintf"),
		// @def fmt.Errorf
		s.ArgExposeUnsafe(ctx, "fmt", "Errorf"),
		// @def fmt.Sscanf
		s.ArgExposeUnsafe(ctx, "fmt", "Sscanf"),
		// @def EnumerationType github.com/xoctopus/x/enumx.Enum
		s.ArgExposeUnsafe(ctx, pkgid, "Enum").WithName("EnumerationType"),
		// @def github.com/xoctopus/x/enumx.Scan
		s.ArgExposeUnsafe(ctx, pkgid, "Scan"),
		// @def github.com/xoctopus/x/enumx.ParseErrorFor
		s.ArgExposeUnsafe(ctx, pkgid, "ParseErrorFor"),

		// @def Type
		s.Arg(ctx, "Type", ident),
		// @def database/sql/driver.Value
		s.ArgT[driver.Value](ctx),
		// @def github.com/xoctopus/x/enumx.DriverValueOffset
		s.ArgT[enumx.DriverValueOffset](ctx),

		// @def UnknownValue
		s.Arg(ctx, "UnknownValue", s.ExposeObject(ctx, e.unknown.Exposer())),
		// @def EnumValues
		s.Arg(ctx, "EnumValues", e.Values(ctx)),
		// @def Values
		s.Arg(ctx, "Values", e.Values(ctx)),
		// @def NameToValueCases
		s.Arg(ctx, "StringToValueCases", e.StringToValueCases(ctx)),
		// @def ValueToDescCases
		s.Arg(ctx, "ValueToTextCases", e.ValueToTextCases(ctx)),
		// @def ValueToNameCases
		s.Arg(ctx, "ValueToStringCases", e.ValueToStringCases(ctx)),
	}

	ss := []s.Snippet{s.Template(bytes.NewReader(template), args...)}
	for _, attr := range e.Attrs() {
		if v := strings.ToLower(attr); v != "text" && v != "string" {
			ss = append(ss, e.Attr(ctx, attr, e.options[attr]))
		}
	}

	switch e.storage {
	case "text", "string", "varchar":
		ss = append(ss, s.Template(bytes.NewReader(templateTxt), args...))
	default:
		ss = append(ss, s.Template(bytes.NewReader(templateInt), args...))
	}

	c.Render(s.Snippets(s.NewLine(1), ss...))
}

func Span() func() time.Duration {
	t := time.Now()
	return func() time.Duration {
		return time.Since(t)
	}
}
