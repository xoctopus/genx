package enumx

// Gender
// @def storage=text
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	// GENDER__MALE 男
	// @attr name=男
	// @attr text=男
	// @attr short=M
	GENDER__MALE
	// GENDER__FEMALE
	// @attr name=女
	// @attr text=女
	// @attr short=F
	// @ignore GENDER__FEMALE has no more description, use key as its Text
	GENDER__FEMALE
	GENDER_INVALID
)
