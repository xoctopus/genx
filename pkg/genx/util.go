package genx

import (
	"strings"

	"github.com/xoctopus/x/slicex"
)

func ParseDirectives(lines ...string) map[string]Directive {
	directives := make(map[string]Directive)
	for _, line := range lines {
		if after, ok := strings.CutPrefix(line, "+genx:"); ok {
			name, parameter, _ := strings.Cut(after, "=")

			if name == "ignore" {
				names := slicex.Mapping(strings.Split(parameter, ","), func(from string) string {
					return strings.TrimSpace(from)
				})

				for _, cmd := range names {
					d := directives[cmd]
					d.enabled = false
					d.name = cmd
					directives[cmd] = d
				}
				continue
			}

			if _, exists := directives[name]; exists {
				continue
			}

			d := Directive{name: name}

			parameter = strings.TrimSpace(parameter)
			if parameter == "false" {
				d.enabled = false
			} else {
				d.enabled = true
				d.parameter = parameter
			}

			directives[name] = d
		}
	}
	return directives
}

type Directive struct {
	name      string
	parameter string
	enabled   bool
}
