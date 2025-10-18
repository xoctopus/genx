package snippet

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
	"sort"
	"strconv"
)

// less only for order map key value
func less(x, y reflect.Value) bool {
	for x.Kind() == reflect.Interface && x.IsValid() {
		x = x.Elem()
	}
	for y.Kind() == reflect.Interface && y.IsValid() {
		y = y.Elem()
	}

	if x.IsValid() && !y.IsValid() {
		return false
	}
	if !x.IsValid() && y.IsValid() {
		return true
	}

	px, py := x.Type().PkgPath(), y.Type().PkgPath()
	if px != py {
		return px < py
	}

	kx, ky := x.Kind(), y.Kind()
	if kx != ky {
		return kx < ky
	}

	sx, sy := fmt.Sprintf("%#v", x.Interface()), fmt.Sprintf("%#v", y.Interface())
	return sx < sy
}

func Value(ctx context.Context, x any) Snippet {
	v, ok := x.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(x)
	}

	switch kind := v.Kind(); kind {
	case reflect.Func:
		return nil // skip unsupported
	case reflect.Pointer:
		if v.IsNil() {
			return Block("nil")
		}
		return Compose(
			Expose(ctx, "github.com/xoctopus/x/ptrx", "Ptr", IdentRT(ctx, v.Elem().Type())),
			Block("("),
			Value(ctx, v.Elem()),
			Block(")"),
		)
	case reflect.Interface:
		if v.IsNil() {
			return Block("nil")
		}
		return Value(ctx, v.Elem())
	case reflect.Map:
		return Compose(
			IdentRT(ctx, v.Type()),
			Block("{"),
			func() Snippet {
				ss := make([]Snippet, 0)

				keys := make([]reflect.Value, 0, v.Len())
				vals := map[reflect.Value]reflect.Value{}
				for _, k := range v.MapKeys() {
					keys = append(keys, k)
					vals[k] = v.MapIndex(k)
				}
				sort.Slice(keys, func(i, j int) bool {
					return less(keys[i], keys[j])
				})
				for i, k := range keys {
					if i > 0 {
						ss = append(ss, Block(", "))
					}
					ss = append(ss, Value(ctx, k))
					ss = append(ss, Block(": "))
					ss = append(ss, Value(ctx, vals[k]))
				}
				return Compose(ss...)
			}(),
			Block("}"),
		)
	case reflect.Slice, reflect.Array:
		return Compose(
			IdentRT(ctx, v.Type()),
			Block("{"),
			func() Snippet {
				ss := make([]Snippet, 0)
				for i := 0; i < v.Len(); i++ {
					if i > 0 {
						ss = append(ss, Block(", "))
					}
					ss = append(ss, Value(ctx, v.Index(i)))
				}
				return Compose(ss...)
			}(),
			Block("}"),
		)
	case reflect.Struct:
		return Compose(
			IdentRT(ctx, v.Type()),
			Block("{"),
			func() Snippet {
				ss := make([]Snippet, 0)
				written := 0
				for i := range v.NumField() {
					fv := v.Field(i)
					ft := v.Type().Field(i)
					if !ast.IsExported(ft.Name) {
						continue
					}
					if written > 0 {
						ss = append(ss, Block(", "))
					}
					ss = append(ss, Block(ft.Name))
					ss = append(ss, Block(": "))
					ss = append(ss, Value(ctx, fv))
					written++
				}
				return Compose(ss...)
			}(),
			Block("}"),
		)
	default: // basics or invalid
		var s Snippet
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s = Block(fmt.Sprintf("%d", v.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			s = Block(fmt.Sprintf("%d", v.Uint()))
		case reflect.Float32:
			s = Block(strconv.FormatFloat(v.Float(), 'f', -1, 32))
		case reflect.Float64:
			s = Block(strconv.FormatFloat(v.Float(), 'f', -1, 64))
		case reflect.String:
			s = Block(strconv.Quote(v.String()))
		case reflect.Complex64:
			s = Block(strconv.FormatComplex(v.Complex(), 'f', -1, 64))
		case reflect.Complex128:
			s = Block(strconv.FormatComplex(v.Complex(), 'f', -1, 128))
		case reflect.Bool:
			s = Block(strconv.FormatBool(v.Bool()))
		default: // reflect.Invalid
			return Block("nil")
		}

		if v.Type().PkgPath() == "" {
			return s
		}
		return Compose(IdentRT(ctx, v.Type()), Block("("), s, Block(")"))
	}
}
