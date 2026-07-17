package testdata

// Gender enum of genders
// +genx:test_genx
// +genx:test_genx_e
// +genx:test_genx_e_a
// +genx:test_genx_t
// +genx:test_genx_ge
// +genx:test_genx_nil
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE          // 男
	GENDER__FEMALE        // 女
)
