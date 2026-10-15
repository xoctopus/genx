package codex

import (
	"go/constant"
	"go/importer"
	"go/types"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

func init() {
	pkg, err := importer.Default().Import("net/http")
	must.NoErrorF(err, "genx/codex: import net/http")

	byName := make(map[string]int)
	byCode := make(map[int]string)

	for _, name := range pkg.Scope().Names() {
		if !strings.HasPrefix(name, "Status") {
			continue
		}
		obj := pkg.Scope().Lookup(name)
		c, ok := obj.(*types.Const)
		if !ok {
			continue
		}
		n, ok := constant.Int64Val(c.Val())
		if !ok || n <= 0 {
			continue
		}
		code := int(n)
		byName[name] = code
		if _, exists := byCode[code]; !exists {
			byCode[code] = name
		}
	}

	gHTTPStatuses = HTTPStatusRegistry{
		byName: byName,
		byCode: byCode,
	}
}

// HTTPStatusRegistry maps net/http status constant names (StatusOK) to numeric codes.
// Populated once from go/types so we stay aligned with the toolchain's net/http.
type HTTPStatusRegistry struct {
	byName map[string]int
	byCode map[int]string
}

var gHTTPStatuses HTTPStatusRegistry

func (r HTTPStatusRegistry) StatusCodeOf(name string) (int, bool) {
	code, ok := r.byName[name]
	return code, ok
}

// ResolveStatusAnnoValue resolves annotation of api
// eg:
//
//	@code api=StatusInternalServerError
//	@code api=http=500
func ResolveStatusAnnoValue(anno string) (string, bool) {
	anno = strings.TrimSpace(anno)
	if code, ok := gHTTPStatuses.StatusCodeOf(anno); ok {
		return strconv.Itoa(code), true
	}
	if s, ok := strings.CutPrefix(anno, "http="); ok {
		s = strings.TrimSpace(s)
		if code, err := strconv.Atoi(s); err == nil && code >= 100 && code <= 999 {
			return strconv.Itoa(code), true
		}
	}
	return "", false
}
