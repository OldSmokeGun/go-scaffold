package types

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// An UnsupportedTypeError is returned by Scan when attempting
// to scan an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "sql value scan: unsupported type: " + e.Type.String()
}

// UnixTimestamp convert unix timestamp in the database to time.Time
// or convert the time.Time to unix timestamp and writes it to the database
type UnixTimestamp struct {
	time.Time
}

func (t UnixTimestamp) Value() (driver.Value, error) {
	return t.Unix(), nil
}

func (t *UnixTimestamp) Scan(src any) (err error) {
	var v int64

	switch i := src.(type) {
	case []byte:
		v, err = strconv.ParseInt(string(i), 10, 64)
		if err != nil {
			return errors.WithStack(err)
		}
	case int64:
		v = i
	default:
		return &UnsupportedTypeError{reflect.TypeOf(i)}
	}

	t.Time = time.Unix(v, 0)
	return nil
}
