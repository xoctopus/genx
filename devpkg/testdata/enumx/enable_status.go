package enumx

// EnableStatus
// +genx:enum
// @enum storage=enum
type EnableStatus int8

const (
	ENABLE_STATUS_UNKNOWN EnableStatus = iota
	// ENABLE_STATUS__ENABLED 开启
	// @enum enum=Y
	ENABLE_STATUS__ENABLED
	// ENABLE_STATUS__DISABLED 关闭
	// @enum enum=N
	ENABLE_STATUS__DISABLED
)
