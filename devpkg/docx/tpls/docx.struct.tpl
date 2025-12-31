@def T
@def FieldDocCases
@def AnonymousDoc
--TypeDoc
func (v *#T#) DocOf(names ...string) ([]string, bool) {
	if len(names) > 0 {
		switch names[0] {
			#FieldDocCases#
		}
		#AnonymousDoc#
		return []string{}, false
	}
	return []string{#TDoc#}, true
}
