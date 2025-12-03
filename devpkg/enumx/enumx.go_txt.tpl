@def Type
@def database/sql/driver.Value
--Value
// Value implements driver.Valuer
func (v #Type#) Value() (#database/sql/driver.Value#, error) {
	return v.String(), nil
}

@def Type
@def fmt.Errorf
--Scan
// Scan implements sql.Scanner
func (v *#Type#) Scan(src any) error {
	var data []byte
	switch x := src.(type) {
	case string:
		data = []byte(x)
	case []byte:
		data = x
	default:
		return #fmt.Errorf#("cannot scan %T value from `%T`", v, x)
	}
	return v.UnmarshalText(data)
}
