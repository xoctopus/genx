package docx

import (
	"strings"

	"github.com/xoctopus/x/slicex"
)

var keywords = map[string]struct{}{
	"XXX":      {},
	"ERROR":    {},
	"BUG":      {},
	"HACK":     {},
	"TODO":     {},
	"PERF":     {},
	"OPTIMIZE": {},
	"REFACTOR": {},
	"NOTE":     {},
	"INFO":     {},
	"NOTICE":   {},
	"WARNING":  {},
	"REVIEW":   {},
}

func Parse(lines []string) *Meta {
	d := &Meta{
		annotations: Annotations{},
		directives:  Directives{},
	}
	if len(lines) == 0 {
		return d
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if parts := strings.Split(line, " "); len(parts) > 0 {
			if _, ok := keywords[strings.ToUpper(parts[0])]; ok {
				continue
			}
		}
		if after, ok := strings.CutPrefix(line, "@"); ok {
			idx := strings.Index(after, " ")
			if idx == 0 || idx == -1 {
				continue
			}
			d.annotations.add(after[0:idx], strings.TrimSpace(after[idx+1:]))
			continue
		}
		if after, ok := strings.CutPrefix(line, "+"); ok {
			if after, ok = strings.CutPrefix(after, "genx:"); ok {
				if !strings.HasPrefix(after, " ") {
					name, parameter, _ := strings.Cut(after, "=")
					d.directives.add(name, parameter)
				}
			}
			continue
		}
		d.description = append(d.description, line)
	}
	return d
}

type Meta struct {
	description []string
	annotations Annotations
	directives  Directives
}

func (m *Meta) Lines() []string {
	return m.description
}

func (m *Meta) Title(prefix string) string {
	lines := m.Lines()
	if len(lines) == 0 {
		return ""
	}
	if len(prefix) == 0 {
		return lines[0]
	}
	if lines[0] == prefix {
		return ""
	}
	return strings.TrimPrefix(lines[0], prefix+" ")
}

func (m *Meta) Description(prefix string, sep string) string {
	return strings.Join(m.Descriptions(prefix), sep)
}

func (m *Meta) Descriptions(prefix string) []string {
	return slicex.FilterMapping(m.description, func(line string) (string, bool) {
		if len(prefix) > 0 {
			if line == prefix {
				return "", false
			}
			line = strings.TrimPrefix(line, prefix+" ")
			return line, len(line) > 0
		}
		return line, true
	})
}

func (m *Meta) Annotations() Annotations {
	return m.annotations
}

func (m *Meta) AnnotationsByName(name string) ([]Annotation, bool) {
	annotations, ok := m.annotations[name]
	return annotations, ok
}

func (m *Meta) Directives() map[string]Directive {
	return m.directives
}

func (m *Meta) Directive(name string) (Directive, bool) {
	d, ok := m.directives[name]
	return d, ok
}

type Annotation struct {
	name string
	text string
}

func (a *Annotation) Name() string {
	return a.name
}

func (a *Annotation) Text() string {
	return a.text
}

func (a *Annotation) Key() string {
	k, _, _ := strings.Cut(a.text, "=")
	return k
}

func (a *Annotation) Value() string {
	_, v, _ := strings.Cut(a.text, "=")
	return v
}

type Annotations map[string][]Annotation

func (as Annotations) add(name, text string) {
	as[name] = append(as[name], Annotation{name: name, text: text})
}

type Directive struct {
	name      string
	parameter string
	enabled   bool
}

func (d *Directive) Name() string {
	return d.name
}

func (d *Directive) Parameter() string {
	return d.parameter
}

func (d *Directive) Enabled() bool {
	return d.enabled
}

type Directives map[string]Directive

func (ds Directives) add(name, parameter string) {
	if len(name) > 0 {
		if name == "ignore" {
			for k := range strings.SplitSeq(parameter, ",") {
				if len(k) > 0 {
					d := ds[k]
					d.name = k
					d.enabled = false
					ds[k] = d
				}
			}
			return
		}
		if _, exists := ds[name]; exists {
			return
		}
		d := Directive{
			name:      name,
			parameter: parameter,
			enabled:   parameter != "false",
		}
		ds[name] = d
	}
}
