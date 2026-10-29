@def CodeType
@def CodeIdentifierCases
-- CodeIdentifier
// CodeIdentifier returns the unique identifier of the code
func (e #CodeType#) CodeIdentifier() string {
	switch e {
	default:
		return ""
	#CodeIdentifierCases#
	}
}
