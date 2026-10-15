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
//
// # Domain in Message()
//
// Default domain is UPPER(package)::UPPER(TypeName). Optional type doc override: `// @code domain=MYDOMAIN`.
// Constant docs use only `api` for HTTP, not domain.
//
// Example:
//
//	// +genx:code
//	type Code int8
//
//	const (
//		CODE_UNDEFINED Code = iota
//		// CODE__ERROR1 error1 message
//		// @code api=StatusBadRequest
//		CODE__ERROR1
//	)
package codex
