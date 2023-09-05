package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go-scaffold/internal/app/repository/schema/mixin"
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
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("phone"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("username").Default("").Comment("用户名"),
		field.String("password").Default("").Comment("密码"),
		field.String("nickname").Default("").Comment("用户名"),
		field.String("phone").Default("").Comment("电话"),
		field.String("salt").Default("").Comment("盐值"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
