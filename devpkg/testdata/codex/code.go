package testdata

// Code internal error code
// +genx:errx_code
type Code int8

const (
	CODE_UNDEFINED Code = iota
	CODE__ERROR1        // error1 message
	CODE__ERROR2        // error2 message
	CODE__ERROR3
	_ // placeholder will be skipped
)
