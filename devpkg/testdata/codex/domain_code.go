package testdata

// DomainCode internal error code has domain name
// +genx:code
// @code domain=DOMAIN_NAME
type DomainCode int8

const (
	DOMAIN_CODE_UNDEFINED  DomainCode = iota
	DOMAIN_CODE__PARSE                // parse failed
	DOMAIN_CODE__HANDLE               // handle failed
	DOMAIN_CODE__PARAMETER            // invalid parameter
)
