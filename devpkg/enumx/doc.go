// Package enumx provides a generator for generating enumeration related implementations.
//
// # How to Integrate
//
//	import _ "github.com/xoctopus/genx/devpkg/enumx"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...)
//
// # Code Definition Conventions
//
//  1. Add `// +genx:enum` comment to your enumeration type.
//  2. Define a custom `int` type and use the `@enum` annotation in its comment
//     to configure the generator (e.g., `// @enum storage=text`).
//  3. Define constants for your type following the naming conventions below.
//
// Assuming your type is named `Status`, its prefix will be `STATUS` (UpperSnakeCase).
//
//   - Unknown Value:
//     Must be named exactly as `<PREFIX>_UNKNOWN`.
//     Example: `STATUS_UNKNOWN`
//
//   - Specific Values:
//     Must be named as `<PREFIX>__<NAME>` (using a double underscore `__`).
//     Example: `STATUS__ENABLED`
//
// The text description is extracted from the constant's comment.
// If no comment is provided, the `<NAME>` is used as the fallback text.
//
// Example:
//
//	// Status
//	// +genx:enum
//	type Status int8
//
//	const (
//		STATUS_UNKNOWN Status = iota
//		// STATUS__ENABLED 关闭
//		STATUS__ENABLED
//		// STATUS__DISABLED 开启
//		STATUS__DISABLED
//	)
package enumx
