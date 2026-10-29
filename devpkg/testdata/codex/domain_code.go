package testdata

// DomainCode internal error code has domain name
// +genx:code
// @code domain=DOMAIN_NAME
type DomainCode int8

const (
	DOMAIN_CODE_UNDEFINED DomainCode = iota
	// DOMAIN_CODE__PARSE parse failed
	// @code hint
	DOMAIN_CODE__PARSE
	DOMAIN_CODE__HANDLE // handle failed
	// DOMAIN_CODE__PARAMETER invalid parameter
	// @code hint=please check the input parameters
	DOMAIN_CODE__PARAMETER
	// DOMAIN_CODE__INTERNAL
	// @code hint
	DOMAIN_CODE__INTERNAL
)
