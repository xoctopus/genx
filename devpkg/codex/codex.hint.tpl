@def CodeType
@def CodeHintCases
-- CodeHint
// CodeHint returns the message which can be presented to users
func (e #CodeType#) CodeHint() string {
	switch e {
	default:
		return ""
	#CodeHintCases#
	}
}
