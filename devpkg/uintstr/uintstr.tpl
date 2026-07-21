@def T
@def strconv.ParseUint
@def strconv.FormatUint
--UintArshaler
// UnmarshalText parse #T#
func (v *#T#) UnmarshalText(text []byte) error {
	s := string(text)
	if len(s) == 0 {
		return nil
	}
	d, err := #strconv.ParseUint#(s, 10, 64)
	if err != nil {
		return err
	}
	*v = #T#(d)
	return nil
}

// MarshalText serializes #T#
func (v #T#) MarshalText() (text []byte, err error) {
	if v == 0 {
		return nil, nil
	}
	return []byte(#strconv.FormatUint#(uint64(v), 10)), nil
}
