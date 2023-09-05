package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go-scaffold/internal/app/repository/schema/mixin"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "permissions",
			Options: "COMMENT='权限表'",
		},
		entsql.WithComments(true),
	}
}

func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key"),
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("key").Unique().MaxLen(128).Comment("权限标识"),
		field.String("name").Default("").MaxLen(128).Comment("权限名称"),
		field.String("desc").Default("").MaxLen(255).Comment("权限描述"),
		field.Int64("parent_id").Default(0).Comment("父级权限 id"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return nil
}
