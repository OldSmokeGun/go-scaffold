package repository

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository/queries"
	"go-scaffold/internal/app/repository/schema"
	"go-scaffold/internal/app/repository/tools"
	igorm "go-scaffold/internal/pkg/gorm"
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
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewPermissionRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *PermissionRepository {
	return &PermissionRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *PermissionRepository) Filter(ctx context.Context, param PermissionFindListParam) ([]*domain.Permission, error) {
	q := gorm.G[schema.Permission](r.db).
		Order(queries.Permission.UpdatedAt.Desc())

	if param.Keyword != "" {
		kw := tools.BuildLikeContains(param.Keyword)
		q = q.
			Where(queries.Permission.Name.Like(kw)).
			Or(queries.Permission.Desc.Like(kw))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.Permission, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *PermissionRepository) FindList(ctx context.Context, idList []int64) ([]*domain.Permission, error) {
	if len(idList) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.ID.In(idList...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	list := make([]*domain.Permission, 0, len(data))
	for i := range data {
		list = append(list, data[i].ToEntity())
	}

	return list, nil
}

func (r *PermissionRepository) FindOne(ctx context.Context, id int64) (*domain.Permission, error) {
	m, err := queries.Query[schema.Permission](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *PermissionRepository) FindOneByKey(ctx context.Context, key string) (*domain.Permission, error) {
	m, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.Key.Eq(key)).
		Take(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *PermissionRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.Permission](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *PermissionRepository) KeyExist(ctx context.Context, key string) (bool, error) {
	n, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.Key.Eq(key)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *PermissionRepository) KeyExistExcludeID(ctx context.Context, key string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.Key.Eq(key)).
		Where(queries.Permission.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *PermissionRepository) HasChild(ctx context.Context, id int64) (bool, error) {
	n, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.ParentID.Eq(id)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *PermissionRepository) Create(ctx context.Context, e domain.Permission) error {
	m := schema.Permission{
		Key:      e.Key,
		Name:     e.Name,
		Desc:     e.Desc,
		ParentID: e.ParentID,
	}

	err := gorm.G[schema.Permission](r.db).
		Create(ctx, &m)

	return handleError(err)
}

func (r *PermissionRepository) Update(ctx context.Context, e domain.Permission) error {
	_, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.ID.Eq(e.ID)).
		Set(
			queries.Permission.Key.Set(e.Key),
			queries.Permission.Name.Set(e.Name),
			queries.Permission.Desc.Set(e.Desc),
			queries.Permission.ParentID.Set(e.ParentID),
		).
		Update(ctx)

	return handleError(err)
}

func (r *PermissionRepository) Delete(ctx context.Context, e domain.Permission) error {
	_, err := r.enforcer.DeletePermission(fmt.Sprintf("%d", e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.Permission](r.db).
		Where(queries.Permission.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}
