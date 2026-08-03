// Package docx provides a generator for extracting and exposing Go source code
// documentation comments at runtime.
//
// # How to Integrate
//
//	import _ "github.com/xoctopus/genx/devpkg/docx"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...)
//
// # Code Definition Conventions
//
// 1. Add `// +genx:doc` to your package-level doc comment.
// 2. Define your types and write standard Go doc comments for them and their fields.
// 3. More comment extraction feature see `github.com/xoctopus/genx/pkg/docx`
//
// The generator extracts documentation based on standard Go comments.
//
// 1. Type Documentation:
// The generator extracts the documentation for named types.
// If the type is a struct, it will also extract documentation for its exported fields.
//
// 2. Field Documentation:
// For struct fields, the documentation is extracted from the field's comment.
// The generator maps the field name to its corresponding comment lines.
//
// 3. Embedded Fields (Anonymous Fields):
// If a struct embeds another type, the generator will attempt to extract documentation for the embedded type as well.
// You can override or prefix the embedded type's documentation by adding a comment to the embedded field. The first line (title) of this comment will be used as a prefix.
//
// Example:
//
//	// Package example
//	// +genx:doc
//	package example
//
//	type UserProfile struct {
//		// Avatar URL
//		Avatar string
//	}
//
//	// User represents a system user.
//	type User struct {
//		// Name is the user's full name.
//		Name string
//		// Age is the user's age in years.
//		Age int
//		// UserProfile user profile
//		UserProfile
//	}
//
// Generated Code
//
//	func (v *UserProfile) DocOf(names ...string) ([]string, bool) {
//		if len(names) > 0 {
//			switch names[0] {
//			case "Avatar":
//				return []string{"URL"}, true
//			}
//			return []string{}, false
//		}
//		return []string{}, true
//	}
//
//	func (v *User) DocOf(names ...string) ([]string, bool) {
//		if len(names) > 0 {
//			switch names[0] {
//			case "Name":
//				return []string{"is the user's full name."}, true
//			case "Age":
//				return []string{"is the user's age in years."}, true
//			}
//			if doc, ok := docx.Of(&v.UserProfile, "user profile", names...); ok {
//				return doc, true
//			}
//			return []string{}, false
//		}
//		return []string{"represents a system user."}, true
//	}
package docx
