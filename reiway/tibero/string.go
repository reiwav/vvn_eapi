package tibero

import (
	"database/sql/driver"
	"fmt"
)

type String string

func (n *String) Scan(value interface{}) error {
	if value == nil {
		*n = ""
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		b := make([]byte, len(v))
		for i, v := range v {
			b[i] = byte(v)
		}
		*n = String(b)
	default:
		*n = String(fmt.Sprintf("%v", v))
	}
	return nil
}

func (n String) Value() (driver.Value, error) {
	return n, nil
}

func (s String) String() string {
	return string(s)
}
