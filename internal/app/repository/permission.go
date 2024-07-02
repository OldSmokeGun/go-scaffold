package repository

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	ient "go-scaffold/internal/pkg/ent"
	"go-scaffold/internal/pkg/ent/ent"
	"go-scaffold/internal/pkg/ent/ent/permission"
)

var _ PermissionRepositoryInterface = (*PermissionRepository)(nil)

type (
	PermissionFindListParam struct {
		Keyword string
	}

	PermissionRepositoryInterface interface {
		Filter(ctx context.Context, param PermissionFindListParam) ([]*domain.Permission, error)
		FindList(ctx context.Context, idList []int64) ([]*domain.Permission, error)
		FindOne(ctx context.Context, id int64) (*domain.Permission, error)
		FindOneByKey(ctx context.Context, key string) (*domain.Permission, error)
		Exist(ctx context.Context, id int64) (bool, error)
		KeyExist(ctx context.Context, key string) (bool, error)
		KeyExistExcludeID(ctx context.Context, key string, excludeID int64) (bool, error)
		HasChild(ctx context.Context, id int64) (bool, error)
		Create(ctx context.Context, e domain.Permission) error
		Update(ctx context.Context, e domain.Permission) error
		Delete(ctx context.Context, e domain.Permission) error
	}
)

type PermissionRepository struct {
	client   *ient.DefaultClient
	enforcer *casbin.Enforcer
}

func NewPermissionRepository(client *ient.DefaultClient, enforcer *casbin.Enforcer) *PermissionRepository {
	return &PermissionRepository{
		client:   client,
		enforcer: enforcer,
	}
}

func (r *PermissionRepository) Filter(ctx context.Context, param PermissionFindListParam) ([]*domain.Permission, error) {
	query := r.client.Permission.Query()

	if param.Keyword != "" {
		query.Where(
			permission.Or(
				permission.NameContains(param.Keyword),
				permission.DescContains(param.Keyword),
			),
		)
	}

	list, err := query.
		Order(ent.Desc(permission.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	entities := make([]*domain.Permission, 0, len(list))
	for _, i := range list {
		entities = append(entities, (&permissionModel{i}).toEntity())
	}

	return entities, nil
}

func (r *PermissionRepository) FindList(ctx context.Context, idList []int64) ([]*domain.Permission, error) {
	data, err := r.client.Permission.Query().
		Where(permission.IDIn(idList...)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	list := make([]*domain.Permission, 0, len(data))
	for _, item := range data {
		list = append(list, (&permissionModel{item}).toEntity())
	}
	return list, nil
}

func (r *PermissionRepository) FindOne(ctx context.Context, id int64) (*domain.Permission, error) {
	m, err := r.client.Permission.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&permissionModel{m}).toEntity(), nil
}

func (r *PermissionRepository) FindOneByKey(ctx context.Context, key string) (*domain.Permission, error) {
	m, err := r.client.Permission.Query().Where(permission.KeyEQ(key)).Only(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&permissionModel{m}).toEntity(), nil
}

func (r *PermissionRepository) Exist(ctx context.Context, id int64) (bool, error) {
	exist, err := r.client.Permission.Query().Where(permission.IDEQ(id)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *PermissionRepository) KeyExist(ctx context.Context, key string) (bool, error) {
	exist, err := r.client.Permission.Query().Where(permission.KeyEQ(key)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *PermissionRepository) KeyExistExcludeID(ctx context.Context, key string, excludeID int64) (bool, error) {
	exist, err := r.client.Permission.Query().Where(
		permission.KeyEQ(key),
		permission.IDNEQ(excludeID),
	).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *PermissionRepository) HasChild(ctx context.Context, id int64) (bool, error) {
	count, err := r.client.Permission.Query().Where(permission.ParentIDEQ(id)).Count(ctx)
	return count > 0, errors.WithStack(handleError(err))
}

func (r *PermissionRepository) Create(ctx context.Context, e domain.Permission) error {
	_, err := r.client.Permission.Create().
		SetKey(e.Key).
		SetName(e.Name).
		SetDesc(e.Desc).
		SetParentID(e.ParentID).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *PermissionRepository) Update(ctx context.Context, e domain.Permission) error {
	_, err := r.client.Permission.
		UpdateOneID(e.ID).
		SetKey(e.Key).
		SetName(e.Name).
		SetDesc(e.Desc).
		SetParentID(e.ParentID).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *PermissionRepository) Delete(ctx context.Context, e domain.Permission) error {
	_, err := r.enforcer.DeletePermission(fmt.Sprintf("%d", e.ID))
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(r.client.Permission.DeleteOneID(e.ID).Exec(ctx))
}

type permissionModel struct {
	*ent.Permission
}

func (m *permissionModel) toEntity() *domain.Permission {
	return &domain.Permission{
		ID:       m.ID,
		Key:      m.Key,
		Name:     m.Name,
		Desc:     m.Desc,
		ParentID: m.ParentID,
	}
}
