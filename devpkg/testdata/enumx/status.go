package enumx

// Status
// @def attr.Gender=Gender
type Status int8

const (
	STATUS_UNKNOWN Status = iota
	// STATUS__ENABLED
	// @attr Gender=GENDER__FEMALE
	STATUS__ENABLED // 关闭
	// STATUS__DISABLED
	// @attr Gender=GENDER__MALE
	STATUS__DISABLED // 开启
	_                // placeholder will be ignored
)
