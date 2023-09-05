package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"go-scaffold/internal/app/repository/schema/types"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			GoType(types.UnixTimestamp{}).
			Immutable().
			Default(func() types.UnixTimestamp {
				return types.UnixTimestamp{Time: time.Now()}
			}),
		field.Time("updated_at").
			GoType(types.UnixTimestamp{}).
			Immutable().
			Default(func() types.UnixTimestamp {
				return types.UnixTimestamp{Time: time.Now()}
			}).
			UpdateDefault(func() types.UnixTimestamp {
				return types.UnixTimestamp{Time: time.Now()}
			}),
	}
}
