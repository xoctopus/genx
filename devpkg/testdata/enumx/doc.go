// Package enumx defines some int/stringer enumerations to generate enumx.Enum
//
// +genx:enum
package enumx

// Weekday not a int/string enum will be skipped
// +genx:enum
type Weekday string

const (
	WEEKDAY__SUN Weekday = "SUN"
	WEEKDAY__MON Weekday = "MON"
	WEEKDAY__TUE Weekday = "TUE"
	WEEKDAY__WED Weekday = "WED"
	WEEKDAY__THU Weekday = "THU"
	WEEKDAY__FRI Weekday = "FRI"
	WEEKDAY__SAT Weekday = "SAT"
)

// GenderExt alias of Gender, not a named type will be skipped
// +genx:enum
type GenderExt = Gender

const (
	GENDER_EXT_UNKNOWN GenderExt = iota
	GENDER_EXT__OTHER
)
