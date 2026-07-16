package genx

import (
	"maps"
	"slices"
	"sort"
	"testing"

	. "github.com/xoctopus/x/testx"
)

func TestParseDirectives(t *testing.T) {
	directives := ParseDirectives(
		"+genx:ignore=a,b,c",
		"+genx:disabled=false",
		"+genx:other=parameter",
		"+genx:a", // exists
	)

	s := slices.Collect(maps.Values(directives))
	sort.Slice(s, func(i, j int) bool {
		return s[i].name < s[j].name
	})
	Expect(t, s, Equal([]Directive{
		{name: "a", enabled: false},
		{name: "b", enabled: false},
		{name: "c", enabled: false},
		{name: "disabled", enabled: false},
		{name: "other", enabled: true, parameter: "parameter"},
	}))
}
