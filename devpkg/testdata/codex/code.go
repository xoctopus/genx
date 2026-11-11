package testdata

// DomainOffset shifts the values of Code to keep them apart from codes of other
// packages. It does not follow `CODE__<NAME>`, so it is not generated as a code.
const DomainOffset = int(100000)

// Code internal error code
// +genx:code
type Code int

const (
	CODE_UNDEFINED Code = iota + Code(DomainOffset)
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
