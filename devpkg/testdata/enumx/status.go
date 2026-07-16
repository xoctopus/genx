package enumx

// Status
// @def attr.Gender=Gender
type Status int8

const (
	STATUS_UNKNOWN Status = iota
	// STATUS__ENABLED 关闭
	// @attr Gender=GENDER__FEMALE
	STATUS__ENABLED
	// STATUS__DISABLED 开启
	// @attr Gender=GENDER__MALE
	STATUS__DISABLED
	_ // placeholder will be ignored
)
