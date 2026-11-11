// Package codex provides a generator for generating error code related implementations.
//
// # How to Integrate Generator
//
//	import _ "github.com/xoctopus/genx/devpkg/codex"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get("code")...)
//
// # Code Definition Conventions
//
//  1. Add `// +genx:code` on the error code type.
//  2. Define constants with `<PREFIX>_UNDEFINED` and `<PREFIX>__<NAME>` (see genx skill generator-codex.spec.md).
//  3. Optional HTTP mapping: on each constant doc, `// @code api=StatusBadRequest` or `// @code api=http=444`.
//     Invalid api values log a warning and omit that constant from StatusCode(). Unannotated codes map to StatusCode() 0; types with no api omit StatusCode().
//  4. Optional hint: on each constant doc, `// @code hint=...` uses the given text, `// @code hint` uses the doc title.
//     Unannotated codes map to CodeHint() ""; types with no hint omit CodeHint().
//  5. Optional offset: start the constants from an offset, such as `CODE_UNDEFINED Code = iota + DomainOffset`,
//     to keep code values apart from other packages. Offsets are managed by the project. The underlying type
//     must hold the offset, and the offset constant must not be named `<PREFIX>__<NAME>`.
//
// # Generated Methods
//
//   - Message() string: `[<domain>:<value>] <doc title or NAME>`.
//   - StatusCode() int: only if any constant has `@code api`.
//   - CodeIdentifier() string: `<KEY_DOMAIN>::<NAME>`, `<KEY_DOMAIN>::UNDEFINED` for the undefined constant.
//   - CodeHint() string: only if any constant has `@code hint`.
//   - CodeValues() []any: all constants except the undefined one, implementing codex.CanBeCode.
//
// # Domain
//
// Optional type doc override: `// @code domain=MYDOMAIN`. Without it:
//
//   - Message() uses UPPER(package)::UPPER(TypeName).
//   - CodeIdentifier() uses UPPER(package), so that keys read as `<PACKAGE>::<NAME>`.
//
// A package is expected to define one code type for its domain. Multiple code
// types in a package are allowed but not recommended: their identifiers share
// the same domain.
//
// Example:
//
//	// +genx:code
//	// @code domain=USER
//	type Code int8
//
//	const (
//		CODE_UNDEFINED Code = iota
//		// CODE__NOT_FOUND user not found
//		// @code api=StatusNotFound
//		// @code hint
//		CODE__NOT_FOUND
//		// CODE__INVALID_NAME invalid name
//		// @code api=StatusBadRequest
//		// @code hint=name must be 8-16 letters or digits
//		CODE__INVALID_NAME
//	)
package codex
