package tibero

import (
	"fmt"
	"reflect"

	"database/sql/driver"
)

type Bool bool

// Scan implements the Scanner interface.
func (n *Bool) Scan(value interface{}) error {
	if value == nil {
		*n = false
		return nil
	}
	fmt.Println(value, reflect.TypeOf(value))
	switch v := value.(type) {
	case int8:
		*n = v > 0
	case int32:
		*n = v > 0
	case int64:
		*n = v > 0
	case float64:
		*n = v > 0
	case float32:
		*n = v > 0
	case bool:
		*n = Bool(v)
	default:
		*n = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n Bool) Value() (driver.Value, error) {
	return n, nil
}
