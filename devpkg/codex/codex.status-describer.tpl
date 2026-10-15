@def CodeType
@def StatusDescriberCases
-- StatusCode
// StatusCode http status describer
func (e #CodeType#) StatusCode() int {
	switch e {
	default:
		return 0
	#StatusDescriberCases#
	}
}
