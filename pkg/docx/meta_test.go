package docx_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/pkg/docx"
)

func ExampleParse() {
	m := docx.Parse([]string{
		" Typename",
		" Typename this line is title ",
		" this is line 1 for descriptions ",
		" this is line 2 for descriptions ",
		" +genx:directive=directive parameters ... ",
		" +genx:ignore=ignored1,ignored2,", // ignored directives
		" +genx:ignored1",                  // redeclared directives
		" + genx:... ",                     // invalid directive will be skipped
		" +genx: ... ",                     // invalid directive will be skipped
		" @attr k1=v1,attr parameters...",
		" TODO prefixed with keywords ",
		" todo keyword is case-sensitive",
		" this is line 3 for descriptions ",
		"", // empty line will be skipped
		" @def annotation content ",
		" @attr k2=v2,attr parameters...",
		" @def ",  // invalid annotation will be skipped
		" @ def ", // invalid annotation will be skipped
	})

	fmt.Println("Description Lines")
	for _, line := range m.Lines() {
		fmt.Println(line)
	}
	fmt.Println()

	fmt.Println("Description String")
	fmt.Println(m.Description("Typename", "\n"))
	fmt.Println()

	fmt.Println("Description String No Prefix")
	fmt.Println(m.Description("", "\n"))
	fmt.Println()

	fmt.Println("Annotations")
	names := slices.Collect(maps.Keys(m.Annotations()))
	slices.Sort(names)
	for _, name := range names {
		fmt.Printf("%s:\n", name)
		annotations, _ := m.AnnotationsByName(name)
		for _, anno := range annotations {
			fmt.Printf("\tname:  %s\n", anno.Name())
			if k := anno.Key(); len(k) > 0 {
				fmt.Printf("\tkey:   %s\n", k)
			}
			if v := anno.Value(); len(v) > 0 {
				fmt.Printf("\tvalue: %s\n", v)
			}
			fmt.Printf("\ttext:  %s\n", anno.Text())
		}
	}
	fmt.Println()

	fmt.Println("Directives")
	names = slices.Collect(maps.Keys(m.Directives()))
	slices.Sort(names)
	for _, key := range names {
		d, _ := m.Directive(key)
		fmt.Printf("%s:\n", key)
		fmt.Printf("\tname:      %s\n", d.Name())
		if p := d.Parameter(); len(p) > 0 {
			fmt.Printf("\tparameter: %s\n", p)
		}
		fmt.Printf("\tenabled:   %t\n", d.Enabled())
	}
	fmt.Println()

	fmt.Println("Empty documents")
	fmt.Println("BEGIN")
	fmt.Println(docx.Parse(nil).Description("", ""))
	fmt.Println("END")

	// Output:
	// Description Lines
	// Typename
	// Typename this line is title
	// this is line 1 for descriptions
	// this is line 2 for descriptions
	// this is line 3 for descriptions
	//
	// Description String
	// this line is title
	// this is line 1 for descriptions
	// this is line 2 for descriptions
	// this is line 3 for descriptions
	//
	// Description String No Prefix
	// Typename
	// Typename this line is title
	// this is line 1 for descriptions
	// this is line 2 for descriptions
	// this is line 3 for descriptions
	//
	// Annotations
	// attr:
	// 	name:  attr
	// 	key:   k1
	// 	value: v1,attr parameters...
	// 	text:  k1=v1,attr parameters...
	// 	name:  attr
	// 	key:   k2
	// 	value: v2,attr parameters...
	// 	text:  k2=v2,attr parameters...
	// def:
	// 	name:  def
	// 	key:   annotation content
	// 	text:  annotation content
	//
	// Directives
	// directive:
	// 	name:      directive
	// 	parameter: directive parameters ...
	// 	enabled:   true
	// ignored1:
	// 	name:      ignored1
	// 	enabled:   false
	// ignored2:
	// 	name:      ignored2
	// 	enabled:   false
	//
	// Empty documents
	// BEGIN
	//
	// END
}

func TestMeta_Title(t *testing.T) {
	m := docx.Parse([]string{"Typename"})
	Expect(t, m.Title("Typename"), Equal(""))

	m = docx.Parse(nil)
	Expect(t, m.Title("Typename"), Equal(""))

	m = docx.Parse([]string{"Typename title"})
	Expect(t, m.Title("Typename"), Equal("title"))
	Expect(t, m.Title(""), Equal("Typename title"))
}
