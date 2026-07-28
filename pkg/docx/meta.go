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

// Parse parses a slice of comment lines into a Meta object.
// It extracts descriptions, annotations, and directives, while skipping lines
// that start with common keywords (like TODO, BUG, etc.).
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

// Meta represents the parsed metadata from a comment block.
type Meta struct {
	description []string
	annotations Annotations
	directives  Directives
}

// Lines returns the raw description lines.
func (m *Meta) Lines() []string {
	return m.description
}

// Title returns the first line of the description.
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

// Description returns the description lines joined by the given separator.
// If prefix is provided, it trims the prefix from each line.
func (m *Meta) Description(prefix string, sep string) string {
	return strings.Join(m.Descriptions(prefix), sep)
}

// Descriptions returns the description lines.
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

// Annotations returns all parsed annotations.
func (m *Meta) Annotations() Annotations {
	return m.annotations
}

// AnnotationsByName returns annotations matching the given name.
func (m *Meta) AnnotationsByName(name string) ([]Annotation, bool) {
	annotations, ok := m.annotations[name]
	return annotations, ok
}

// Directives returns all parsed directives.
func (m *Meta) Directives() map[string]Directive {
	return m.directives
}

// Directive returns a directive by its name.
func (m *Meta) Directive(name string) (Directive, bool) {
	d, ok := m.directives[name]
	return d, ok
}

// Annotation represents a parsed annotation (e.g., @name text).
type Annotation struct {
	name string
	text string
}

// Name returns the name of the annotation.
func (a *Annotation) Name() string {
	return a.name
}

// Text returns the raw text following the annotation name.
func (a *Annotation) Text() string {
	return a.text
}

// Key returns the key part if the annotation text is in key=value format.
func (a *Annotation) Key() string {
	k, _, _ := strings.Cut(a.text, "=")
	return k
}

// Value returns the value part if the annotation text is in key=value format.
func (a *Annotation) Value() string {
	_, v, _ := strings.Cut(a.text, "=")
	return v
}

// Annotations represents a collection of annotations grouped by name.
type Annotations map[string][]Annotation

func (as Annotations) add(name, text string) {
	as[name] = append(as[name], Annotation{name: name, text: text})
}

// Directive represents a parsed directive (e.g., +genx:name=parameter).
type Directive struct {
	name      string
	parameter string
	enabled   bool
}

// Name returns the name of the directive.
func (d *Directive) Name() string {
	return d.name
}

// Parameter returns the parameter string of the directive.
func (d *Directive) Parameter() string {
	return d.parameter
}

// Enabled returns true if the directive is enabled.
// Directives are enabled by default unless the parameter is "false" or it is explicitly ignored.
func (d *Directive) Enabled() bool {
	return d.enabled
}

// Directives represents a collection of directives keyed by name.
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
