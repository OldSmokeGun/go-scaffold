package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go-scaffold/internal/app/repository/schema/mixin"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

func (Product) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "products",
			Options: "COMMENT='产品表'",
		},
		entsql.WithComments(true),
	}
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.SoftDeleteMixin{},
	}
}

func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable(),
		field.String("name").Default("").Comment("名称"),
		field.String("desc").Default("").Comment("描述"),
		field.Int("price").Default(0).Comment("价格"),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return nil
}
