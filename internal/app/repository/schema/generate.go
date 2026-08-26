package schema

import (
	"gorm.io/cli/gorm/field"
	"gorm.io/cli/gorm/genconfig"
	"gorm.io/plugin/soft_delete"
)

//go:generate go run gorm.io/cli/gorm gen -i . -o ../queries

var _ = genconfig.Config{
	OutPath: "../queries",
	FieldTypeMap: map[any]any{
		soft_delete.DeletedAt(0): field.Number[int64]{},
	},
	FieldNameMap: map[string]any{
		"softDelete": field.Number[int64]{},
	},
	IncludeInterfaces: []any{"Query*"},
	IncludeStructs: []any{
		User{},
		Role{},
		Permission{},
		Product{},
	},
}
