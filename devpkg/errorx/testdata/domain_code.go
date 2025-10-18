package testdata

// DomainCode internal error code has domain name
// +genx:code_error
// @def DOMAIN_NAME
type DomainCode int8

const (
	DOMAIN_CODE_UNDEFINED  DomainCode = iota
	DOMAIN_CODE__PARSE                // parse failed
	DOMAIN_CODE__HANDLE               // handle failed
	DOMAIN_CODE__PARAMETER            // invalid parameter
)
