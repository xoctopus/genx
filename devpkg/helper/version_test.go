package helper

import (
	"testing"

	"github.com/xoctopus/x/testx/bdd"
)

func TestVersionFor(t *testing.T) {
	bdd.From(t).Given("empty target to get main module version", func(t bdd.T) {
		version := VersionFor("")
		t.Then(
			"assert main module version",
			bdd.Equal(version, "(devel)"),
		)
	})
	bdd.From(t).Given("other target", func(t bdd.T) {
		version := VersionFor("any other")
		t.Then(
			"assert target version",
			bdd.Equal(version, "devel"),
		)
	})
}
