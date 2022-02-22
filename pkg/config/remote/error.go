package remote

import "fmt"

type FieldMissingError struct {
	Field string
}

func (e FieldMissingError) Error() string {
	return fmt.Sprintf("remote config optional field %s is missing", e.Field)
}

type FieldTypeConvertError struct {
	Field string
	Type  string
}

func (e FieldTypeConvertError) Error() string {
	return fmt.Sprintf("remote config optional field %s can't convert to %s type", e.Field, e.Type)
}
