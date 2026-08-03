// Package codex provides a generator for generating error code related implementations.
//
// # How to Integrate Generator
//
//	import _ "github.com/xoctopus/genx/devpkg/codex"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...)
//
// # Code Definition Conventions
//
//  1. Add `// +genx:code` comment to your error code type.
//  2. Define a custom `int` type and use the `@code` annotation in its comment
//     to configure the domain of error (e.g., `// @code domain="OUR_DOMAIN`).
//  3. Define constants for your type following the naming conventions below.
//
// Assuming your type is named `ErrorCode`, its prefix will be `ERROR_CODE` (UpperSnakeCase).
//
//   - Undefined Error:
//     Must be named exactly as `<PREFIX>_UNDEFINED`.
//     Example: `ERROR_CODE_UNDEFINED`
//
//   - Specific Errors:
//     Must be named as `<PREFIX>__<ERROR_NAME>` (using a double underscore `__`).
//     Example: `ERROR_CODE__NOT_FOUND`
//
// The error message is extracted from the constant's comment. If no comment is
// provided, the `<ERROR_NAME>` is used as the fallback message.
//
// Example:
//
//	// Code internal error code
//	// +genx:code
//	type Code int8
//
//	const (
//		CODE_UNDEFINED Code = iota
//		CODE__ERROR1        // error1 message
//		CODE__ERROR2        // error2 message
//		CODE__ERROR3
//	)
package codex
