package ignored

// Gender enum of genders
// +genx:test_genx
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE          // 男
	GENDER__FEMALE        // 女
)
