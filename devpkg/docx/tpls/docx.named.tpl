@def T
@def TDoc
--TypeDoc
func (v *#T#) DocOf(names ...string) ([]string, bool) {
	if len(names) > 0 {
		return []string{}, false
	}
	return []string{#TDoc#}, true
}
