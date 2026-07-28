// Package docx provides utilities for parsing and extracting metadata from Go
// source code comments.
//
// # Comment Structure
//
// A parsed comment block consists of three main parts:
//
//  1. Directives:
//     Lines starting with `+genx:`. These are used to control the behavior of
//     the generator.
//     Format: `// +genx:name=parameter`
//
//  2. Annotations:
//     Lines starting with `@`. These are used to attach arbitrary kv pairs or
//     text to the element.
//     Format: `// @name text` or `// @name key=value`
//
//  3. Description (Title + Tails):
//     All other lines that are not directives, annotations, or prefixed with
//     keywords (like TODO, BUG).
//     The first line of the description is considered the "Title".
//     The remaining lines are considered the "Tails" (or detailed description).
//     Empty lines are skipped.
package docx
