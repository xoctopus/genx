
@def CodeType
@def CodeMessageCases
@def UnknownCodeFormat
-- CodeMessage
func (e #CodeType#) Message() string {
	switch e {
	default:
		return fmt.Sprintf(#UnknownCodeFormat#, e)
	#CodeMessageCases#
	}
}
