@def Type
@def database/sql/driver.Value
@def github.com/xoctopus/x/enumx.DriverValueOffset
--Value
// Value implements driver.Valuer
func (v #Type#) Value() (#database/sql/driver.Value#, error) {
	offset := 0
	if drv, ok := any(v).(#github.com/xoctopus/x/enumx.DriverValueOffset#); ok {
		offset = drv.Offset()
	}
	return int64(v) + int64(offset), nil
}

@def Type
@def github.com/xoctopus/x/enumx.DriverValueOffset
@def github.com/xoctopus/x/enumx.Scan
--Scan
// Scan implements sql.Scanner
func (v *#Type#) Scan(src any) error {
	offset := 0
	if offsetter, ok := any(v).(#github.com/xoctopus/x/enumx.DriverValueOffset#); ok {
		offset = offsetter.Offset()
	}
	i, err := #github.com/xoctopus/x/enumx.Scan#(src, offset)
	if err != nil {
		return err
	}
	*v = #Type#(i)
	return nil
}
