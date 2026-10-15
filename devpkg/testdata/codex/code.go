package testdata

// Code internal error code
// +genx:code
type Code int8

const (
	CODE_UNDEFINED Code = iota
	// CODE__ERROR1 error1 message
	// @code api=StatusBadRequest
	CODE__ERROR1
	// CODE__ERROR2 error2 message
	// @code api=http=444
	CODE__ERROR2
	// CODE__ERROR3
	// @other x=y
	CODE__ERROR3
	_ // placeholder will be skipped
	// CODE__ERROR4
	// @code api=StatusInvalid
	CODE__ERROR4
)
