package codex_test

import (
	"net/http"
	"strconv"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/devpkg/codex"
)

func TestHTTPStatusRegistry(t *testing.T) {
	code, ok := codex.ResolveStatusAnnoValue("StatusInternalServerError")
	Expect(t, ok, BeTrue())
	Expect(t, code, Equal(strconv.Itoa(http.StatusInternalServerError)))

	code, ok = codex.ResolveStatusAnnoValue("http=444")
	Expect(t, ok, BeTrue())
	Expect(t, code, Equal("444"))

	code, ok = codex.ResolveStatusAnnoValue("StatusInvalid")
	Expect(t, ok, BeFalse())

	code, ok = codex.ResolveStatusAnnoValue("444")
	Expect(t, ok, BeFalse())
}
