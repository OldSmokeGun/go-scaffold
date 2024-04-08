package mixin

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"go-scaffold/internal/app/repository/schema/types"
	gen "go-scaffold/internal/pkg/ent/ent"
	"go-scaffold/internal/pkg/ent/ent/hook"
	"go-scaffold/internal/pkg/ent/ent/intercept"
)

// UnexpectedMutationTypeError mutation does not implement the specified interface
type UnexpectedMutationTypeError struct {
	mutation ent.Mutation
}

func (e *UnexpectedMutationTypeError) Error() string {
	return fmt.Sprintf("unexpected mutation type %T", e.mutation)
}

type SoftDeleteMixin struct {
	mixin.Schema
}

// Fields of the SoftDeleteMixin.
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			GoType(types.UnixTimestamp{}).
			Optional(),
	}
}

// Interceptors of the SoftDeleteMixin.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if isSkipSoftDelete(ctx) {
				return nil
			}
			d.P(q)
			return nil
		}),
	}
}

// Hooks of the SoftDeleteMixin.
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					// skip soft-delete, means delete the entity permanently.
					if isSkipSoftDelete(ctx) {
						return next.Mutate(ctx, m)
					}

					mx, ok := m.(interface {
						SetOp(ent.Op)
						Client() *gen.Client
						SetDeletedAt(types.UnixTimestamp)
						WhereP(...func(*sql.Selector))
					})
					if !ok {
						return nil, &UnexpectedMutationTypeError{m}
					}
					d.P(mx)
					mx.SetOp(ent.OpUpdate)
					mx.SetDeletedAt(types.UnixTimestamp{Time: time.Now()})
					return mx.Client().Mutate(ctx, m)
				})
			},
			ent.OpDeleteOne|ent.OpDelete,
		),
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldEQ(d.Fields()[0].Descriptor().Name, 0),
	)
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, softDeleteKey{}, true)
}

func isSkipSoftDelete(ctx context.Context) bool {
	skip, _ := ctx.Value(softDeleteKey{}).(bool)
	return skip
}
