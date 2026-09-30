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

var _ SystemPermissionRepositoryInterface = (*SystemPermissionRepository)(nil)

type (
	SystemPermissionFindListParam struct {
		Keyword string
	}

	SystemPermissionRepositoryInterface interface {
		Filter(ctx context.Context, param SystemPermissionFindListParam) ([]*domain.SystemPermission, error)
		FindList(ctx context.Context, idList []int64) ([]*domain.SystemPermission, error)
		FindOne(ctx context.Context, id int64) (*domain.SystemPermission, error)
		FindOneByKey(ctx context.Context, key string) (*domain.SystemPermission, error)
		Exist(ctx context.Context, id int64) (bool, error)
		KeyExist(ctx context.Context, key string) (bool, error)
		KeyExistExcludeID(ctx context.Context, key string, excludeID int64) (bool, error)
		HasChild(ctx context.Context, id int64) (bool, error)
		Create(ctx context.Context, e domain.SystemPermission) error
		Update(ctx context.Context, e domain.SystemPermission) error
		Delete(ctx context.Context, e domain.SystemPermission) error
	}
)

type SystemPermissionRepository struct {
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewSystemPermissionRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *SystemPermissionRepository {
	return &SystemPermissionRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *SystemPermissionRepository) Filter(ctx context.Context, param SystemPermissionFindListParam) ([]*domain.SystemPermission, error) {
	q := gorm.G[schema.SystemPermission](r.db).
		Order(queries.SystemPermission.UpdatedAt.Desc())

	if param.Keyword != "" {
		kw := tools.BuildLikeContains(param.Keyword)
		q = q.
			Where(queries.SystemPermission.Name.Like(kw)).
			Or(queries.SystemPermission.Desc.Like(kw))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.SystemPermission, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *SystemPermissionRepository) FindList(ctx context.Context, idList []int64) ([]*domain.SystemPermission, error) {
	if len(idList) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ID.In(idList...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	list := make([]*domain.SystemPermission, 0, len(data))
	for i := range data {
		list = append(list, data[i].ToEntity())
	}

	return list, nil
}

func (r *SystemPermissionRepository) FindOne(ctx context.Context, id int64) (*domain.SystemPermission, error) {
	m, err := queries.Query[schema.SystemPermission](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemPermissionRepository) FindOneByKey(ctx context.Context, key string) (*domain.SystemPermission, error) {
	m, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.Key.Eq(key)).
		Take(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemPermissionRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.SystemPermission](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *SystemPermissionRepository) KeyExist(ctx context.Context, key string) (bool, error) {
	n, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.Key.Eq(key)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemPermissionRepository) KeyExistExcludeID(ctx context.Context, key string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.Key.Eq(key)).
		Where(queries.SystemPermission.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemPermissionRepository) HasChild(ctx context.Context, id int64) (bool, error) {
	n, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ParentID.Eq(id)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemPermissionRepository) Create(ctx context.Context, e domain.SystemPermission) error {
	m := schema.SystemPermission{
		Key:      e.Key,
		Name:     e.Name,
		Desc:     e.Desc,
		ParentID: e.ParentID,
	}

	err := gorm.G[schema.SystemPermission](r.db).
		Create(ctx, &m)

	return handleError(err)
}

func (r *SystemPermissionRepository) Update(ctx context.Context, e domain.SystemPermission) error {
	_, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ID.Eq(e.ID)).
		Set(
			queries.SystemPermission.Key.Set(e.Key),
			queries.SystemPermission.Name.Set(e.Name),
			queries.SystemPermission.Desc.Set(e.Desc),
			queries.SystemPermission.ParentID.Set(e.ParentID),
		).
		Update(ctx)

	return handleError(err)
}

func (r *SystemPermissionRepository) Delete(ctx context.Context, e domain.SystemPermission) error {
	_, err := r.enforcer.DeletePermission(fmt.Sprintf("%d", e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}
