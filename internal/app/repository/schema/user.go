package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "users",
			Options: "COMMENT='用户表'",
		},
		entsql.WithComments(true),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("").Comment("名称"),
		field.Int("age").Default(0).Positive().Comment("年龄"),
		field.String("phone").Default("").Comment("电话"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
