// Package snippet provides a flexible and composable AST-aware code generation
// framework.
//
// Unlike simple string concatenation, this package is designed to work closely
// with Go's syntax tree (AST) and type system. It automatically handles complex
// generation tasks like tracking package imports, resolving type identifiers,
// and safely exposing external objects.
//
// The core concept is the `Snippet` interface, which represents a fragment
// of code that can be generated dynamically. Snippets can be composed,
// nested, and evaluated lazily.
//
// # Core Concepts
//
//   - Snippet: The fundamental interface representing a piece of code.
//     It provides `Fragments` to yield strings iteratively, and `IsNil`
//     to check if it's empty.
//   - Composition: Snippets can be combined using `Compose`, `Snippets`
//     (with separators), or `Strings`.
//
// # Available Snippet Types
//
//  1. Basic Blocks (`block.go`):
//     - `Block(string)`: A raw string block.
//     - `BlockF(format, args...)`: A formatted string block.
//     - `BlockRaw(string)`: A string block that is automatically quoted.
//
//  2. Identifiers (`ident.go`):
//     Used to generate Go type identifiers, automatically handling imports.
//     - `IdentFor[T]()`: Generates the identifier for a generic type `T`.
//     - `IdentOf(v)`: Generates the identifier for the type of value `v`.
//     - `Ident(...)`: Generates identifiers from typx, reflect, or types.
//
//  3. Values (`value.go`):
//     - `Value(v)`: Generates the Go code representation of a value
//     (e.g., structs, maps, slices, primitives).
//
//  4. Code Exposing (`expose.go`):
//     Used to reference exported objects (functions, variables, constants,
//     types) from other packages.
//     - `Expose(path, name, targs...)`: References an object by path and name.
//     - `ExposeObject(types.Object, targs...)`: References a `types.Object`.
//
//  5. Templates (`template.go`):
//     A text-based template engine that replaces macros (e.g., `#MacroName#`)
//     with other Snippets.
//     - `Template(io.Reader, args...)`: Parses a template file.
//     - `Arg(name, snippet)`: Defines arguments for the template.
//
//  6. Comments & Directives (`comment.go`):
//     - `Document(name, lines...)`: Generates a standard GoDoc comment block.
//     - `InlineComment(line)`: Generates an inline comment.
//     - `Directive(name, args...)`: Generates a Go compiler directive.
//
//  7. Imports (`imports.go`):
//     - `Imports()`: Automatically collects and formats all package imports
//     tracked during snippet evaluation.
//
//  8. Utilities (`separator.go`, `poster.go`):
//     - `Indent(n)` / `NewLine(n)`: Generates whitespace separators.
//     - `Poster(pkg, gen, ver)`: Generates the `// Code generated` header.
package snippet
