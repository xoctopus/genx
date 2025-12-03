package testdata

// Code2 internal error code
// +genx:errx_code
// alias will be skipped
type Code2 = int

const (
	CODE_2_UNKNOWN Code2 = iota
	CODE_2__ERROR1
)

// Code3 internal error code
// +genx:errx_code
// string const will be skipped
type Code3 string

const (
	CODE_3_UNKNOWN Code3 = "undefined"
	CODE_3__ERROR1 Code3 = "error1"
)
