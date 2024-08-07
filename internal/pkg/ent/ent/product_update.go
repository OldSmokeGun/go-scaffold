// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"go-scaffold/internal/app/repository/schema/types"
	"go-scaffold/internal/pkg/ent/ent/predicate"
	"go-scaffold/internal/pkg/ent/ent/product"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ProductUpdate is the builder for updating Product entities.
type ProductUpdate struct {
	config
	hooks     []Hook
	mutation  *ProductMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ProductUpdate builder.
func (pu *ProductUpdate) Where(ps ...predicate.Product) *ProductUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetDeletedAt sets the "deleted_at" field.
func (pu *ProductUpdate) SetDeletedAt(tt types.UnixTimestamp) *ProductUpdate {
	pu.mutation.SetDeletedAt(tt)
	return pu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableDeletedAt(tt *types.UnixTimestamp) *ProductUpdate {
	if tt != nil {
		pu.SetDeletedAt(*tt)
	}
	return pu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (pu *ProductUpdate) ClearDeletedAt() *ProductUpdate {
	pu.mutation.ClearDeletedAt()
	return pu
}

// SetName sets the "name" field.
func (pu *ProductUpdate) SetName(s string) *ProductUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableName(s *string) *ProductUpdate {
	if s != nil {
		pu.SetName(*s)
	}
	return pu
}

// SetDesc sets the "desc" field.
func (pu *ProductUpdate) SetDesc(s string) *ProductUpdate {
	pu.mutation.SetDesc(s)
	return pu
}

// SetNillableDesc sets the "desc" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableDesc(s *string) *ProductUpdate {
	if s != nil {
		pu.SetDesc(*s)
	}
	return pu
}

// SetPrice sets the "price" field.
func (pu *ProductUpdate) SetPrice(i int) *ProductUpdate {
	pu.mutation.ResetPrice()
	pu.mutation.SetPrice(i)
	return pu
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (pu *ProductUpdate) SetNillablePrice(i *int) *ProductUpdate {
	if i != nil {
		pu.SetPrice(*i)
	}
	return pu
}

// AddPrice adds i to the "price" field.
func (pu *ProductUpdate) AddPrice(i int) *ProductUpdate {
	pu.mutation.AddPrice(i)
	return pu
}

// Mutation returns the ProductMutation object of the builder.
func (pu *ProductUpdate) Mutation() *ProductMutation {
	return pu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProductUpdate) Save(ctx context.Context) (int, error) {
	if err := pu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProductUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProductUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProductUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *ProductUpdate) defaults() error {
	if _, ok := pu.mutation.UpdatedAt(); !ok {
		if product.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized product.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := product.UpdateDefaultUpdatedAt()
		pu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pu *ProductUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ProductUpdate {
	pu.modifiers = append(pu.modifiers, modifiers...)
	return pu
}

func (pu *ProductUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(product.Table, product.Columns, sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt64))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdatedAt(); ok {
		_spec.SetField(product.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pu.mutation.DeletedAt(); ok {
		_spec.SetField(product.FieldDeletedAt, field.TypeTime, value)
	}
	if pu.mutation.DeletedAtCleared() {
		_spec.ClearField(product.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(product.FieldName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Desc(); ok {
		_spec.SetField(product.FieldDesc, field.TypeString, value)
	}
	if value, ok := pu.mutation.Price(); ok {
		_spec.SetField(product.FieldPrice, field.TypeInt, value)
	}
	if value, ok := pu.mutation.AddedPrice(); ok {
		_spec.AddField(product.FieldPrice, field.TypeInt, value)
	}
	_spec.AddModifiers(pu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{product.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ProductUpdateOne is the builder for updating a single Product entity.
type ProductUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ProductMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetDeletedAt sets the "deleted_at" field.
func (puo *ProductUpdateOne) SetDeletedAt(tt types.UnixTimestamp) *ProductUpdateOne {
	puo.mutation.SetDeletedAt(tt)
	return puo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableDeletedAt(tt *types.UnixTimestamp) *ProductUpdateOne {
	if tt != nil {
		puo.SetDeletedAt(*tt)
	}
	return puo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (puo *ProductUpdateOne) ClearDeletedAt() *ProductUpdateOne {
	puo.mutation.ClearDeletedAt()
	return puo
}

// SetName sets the "name" field.
func (puo *ProductUpdateOne) SetName(s string) *ProductUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableName(s *string) *ProductUpdateOne {
	if s != nil {
		puo.SetName(*s)
	}
	return puo
}

// SetDesc sets the "desc" field.
func (puo *ProductUpdateOne) SetDesc(s string) *ProductUpdateOne {
	puo.mutation.SetDesc(s)
	return puo
}

// SetNillableDesc sets the "desc" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableDesc(s *string) *ProductUpdateOne {
	if s != nil {
		puo.SetDesc(*s)
	}
	return puo
}

// SetPrice sets the "price" field.
func (puo *ProductUpdateOne) SetPrice(i int) *ProductUpdateOne {
	puo.mutation.ResetPrice()
	puo.mutation.SetPrice(i)
	return puo
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillablePrice(i *int) *ProductUpdateOne {
	if i != nil {
		puo.SetPrice(*i)
	}
	return puo
}

// AddPrice adds i to the "price" field.
func (puo *ProductUpdateOne) AddPrice(i int) *ProductUpdateOne {
	puo.mutation.AddPrice(i)
	return puo
}

// Mutation returns the ProductMutation object of the builder.
func (puo *ProductUpdateOne) Mutation() *ProductMutation {
	return puo.mutation
}

// Where appends a list predicates to the ProductUpdate builder.
func (puo *ProductUpdateOne) Where(ps ...predicate.Product) *ProductUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProductUpdateOne) Select(field string, fields ...string) *ProductUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Product entity.
func (puo *ProductUpdateOne) Save(ctx context.Context) (*Product, error) {
	if err := puo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProductUpdateOne) SaveX(ctx context.Context) *Product {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProductUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProductUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *ProductUpdateOne) defaults() error {
	if _, ok := puo.mutation.UpdatedAt(); !ok {
		if product.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized product.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := product.UpdateDefaultUpdatedAt()
		puo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (puo *ProductUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ProductUpdateOne {
	puo.modifiers = append(puo.modifiers, modifiers...)
	return puo
}

func (puo *ProductUpdateOne) sqlSave(ctx context.Context) (_node *Product, err error) {
	_spec := sqlgraph.NewUpdateSpec(product.Table, product.Columns, sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt64))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Product.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, product.FieldID)
		for _, f := range fields {
			if !product.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != product.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.UpdatedAt(); ok {
		_spec.SetField(product.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := puo.mutation.DeletedAt(); ok {
		_spec.SetField(product.FieldDeletedAt, field.TypeTime, value)
	}
	if puo.mutation.DeletedAtCleared() {
		_spec.ClearField(product.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(product.FieldName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Desc(); ok {
		_spec.SetField(product.FieldDesc, field.TypeString, value)
	}
	if value, ok := puo.mutation.Price(); ok {
		_spec.SetField(product.FieldPrice, field.TypeInt, value)
	}
	if value, ok := puo.mutation.AddedPrice(); ok {
		_spec.AddField(product.FieldPrice, field.TypeInt, value)
	}
	_spec.AddModifiers(puo.modifiers...)
	_node = &Product{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{product.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
