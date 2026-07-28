// Package uintstr provides a generator that implements the encoding.TextMarshaler
// and encoding.TextUnmarshaler interfaces for unsigned integer types.
//
// When applied to a named unsigned integer type (e.g., type MyID uint64),
// this generator will automatically create `MarshalText` and `UnmarshalText`
// methods, allowing the type to be serialized to and from its string
// representation (base 10) in formats like JSON, XML, or plain text.
//
// # How to Integrate
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...))
//
// # Code Definition Conventions
//
// 1. Define a named type based on an unsigned integer
// (uint, uint8, uint16, uint32, uint64).
// 2. Add the `// +genx:uintstr` directive to the type's documentation comment.
//
// Example:
//
//	// +genx:uintstr
//	type UserID uint64
//
// Generated Code:
//
//	// UnmarshalText parse UserID
//	func (v *UserID) UnmarshalText(text []byte) error {
//		s := string(text)
//		if len(s) == 0 {
//			return nil
//		}
//		d, err := strconv.ParseUint(s, 10, 64)
//		if err != nil {
//			return err
//		}
//		*v = UserID(d)
//		return nil
//	}
//
//	// MarshalText serializes UserID
//	func (v UserID) MarshalText() (text []byte, err error) {
//		if v == 0 {
//			return nil, nil
//		}
//		return []byte(strconv.FormatUint(uint64(v), 10)), nil
//	}
package uintstr
