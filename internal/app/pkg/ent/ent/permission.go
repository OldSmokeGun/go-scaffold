// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"go-scaffold/internal/app/pkg/ent/ent/permission"
	"go-scaffold/internal/app/repository/schema/types"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Permission is the model entity for the Permission schema.
type Permission struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt types.UnixTimestamp `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt types.UnixTimestamp `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt types.UnixTimestamp `json:"deleted_at,omitempty"`
	// 权限标识
	Key string `json:"key,omitempty"`
	// 权限名称
	Name string `json:"name,omitempty"`
	// 权限描述
	Desc string `json:"desc,omitempty"`
	// 父级权限 id
	ParentID     int64 `json:"parent_id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Permission) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case permission.FieldID, permission.FieldParentID:
			values[i] = new(sql.NullInt64)
		case permission.FieldKey, permission.FieldName, permission.FieldDesc:
			values[i] = new(sql.NullString)
		case permission.FieldCreatedAt, permission.FieldUpdatedAt, permission.FieldDeletedAt:
			values[i] = new(types.UnixTimestamp)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Permission fields.
func (pe *Permission) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case permission.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pe.ID = int64(value.Int64)
		case permission.FieldCreatedAt:
			if value, ok := values[i].(*types.UnixTimestamp); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value != nil {
				pe.CreatedAt = *value
			}
		case permission.FieldUpdatedAt:
			if value, ok := values[i].(*types.UnixTimestamp); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value != nil {
				pe.UpdatedAt = *value
			}
		case permission.FieldDeletedAt:
			if value, ok := values[i].(*types.UnixTimestamp); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value != nil {
				pe.DeletedAt = *value
			}
		case permission.FieldKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key", values[i])
			} else if value.Valid {
				pe.Key = value.String
			}
		case permission.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pe.Name = value.String
			}
		case permission.FieldDesc:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field desc", values[i])
			} else if value.Valid {
				pe.Desc = value.String
			}
		case permission.FieldParentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field parent_id", values[i])
			} else if value.Valid {
				pe.ParentID = value.Int64
			}
		default:
			pe.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Permission.
// This includes values selected through modifiers, order, etc.
func (pe *Permission) Value(name string) (ent.Value, error) {
	return pe.selectValues.Get(name)
}

// Update returns a builder for updating this Permission.
// Note that you need to call Permission.Unwrap() before calling this method if this Permission
// was returned from a transaction, and the transaction was committed or rolled back.
func (pe *Permission) Update() *PermissionUpdateOne {
	return NewPermissionClient(pe.config).UpdateOne(pe)
}

// Unwrap unwraps the Permission entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pe *Permission) Unwrap() *Permission {
	_tx, ok := pe.config.driver.(*txDriver)
	if !ok {
		panic("ent: Permission is not a transactional entity")
	}
	pe.config.driver = _tx.drv
	return pe
}

// String implements the fmt.Stringer.
func (pe *Permission) String() string {
	var builder strings.Builder
	builder.WriteString("Permission(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pe.ID))
	builder.WriteString("created_at=")
	builder.WriteString(fmt.Sprintf("%v", pe.CreatedAt))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(fmt.Sprintf("%v", pe.UpdatedAt))
	builder.WriteString(", ")
	builder.WriteString("deleted_at=")
	builder.WriteString(fmt.Sprintf("%v", pe.DeletedAt))
	builder.WriteString(", ")
	builder.WriteString("key=")
	builder.WriteString(pe.Key)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(pe.Name)
	builder.WriteString(", ")
	builder.WriteString("desc=")
	builder.WriteString(pe.Desc)
	builder.WriteString(", ")
	builder.WriteString("parent_id=")
	builder.WriteString(fmt.Sprintf("%v", pe.ParentID))
	builder.WriteByte(')')
	return builder.String()
}

// Permissions is a parsable slice of Permission.
type Permissions []*Permission