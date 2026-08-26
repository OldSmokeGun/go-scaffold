package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository/queries"
	"go-scaffold/internal/app/repository/schema"
	"go-scaffold/internal/app/repository/tools"
	igorm "go-scaffold/internal/pkg/gorm"
)

var _ RoleRepositoryInterface = (*RoleRepository)(nil)

type (
	RoleFindListParam struct {
		Keyword string
	}

	RoleRepositoryInterface interface {
		Filter(ctx context.Context, param RoleFindListParam) ([]*domain.Role, error)
		FindList(ctx context.Context, idList []int64) ([]*domain.Role, error)
		FindOne(ctx context.Context, id int64) (*domain.Role, error)
		Exist(ctx context.Context, id int64) (bool, error)
		NameExist(ctx context.Context, name string) (bool, error)
		NameExistExcludeID(ctx context.Context, name string, excludeID int64) (bool, error)
		Create(ctx context.Context, e domain.Role) error
		Update(ctx context.Context, e domain.Role) error
		Delete(ctx context.Context, e domain.Role) error
		GrantPermissions(ctx context.Context, role int64, permissions []int64) error
		GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error)
	}
)

type RoleRepository struct {
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewRoleRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *RoleRepository {
	return &RoleRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *RoleRepository) Filter(ctx context.Context, param RoleFindListParam) ([]*domain.Role, error) {
	q := gorm.G[schema.Role](r.db).
		Order(queries.Role.UpdatedAt.Desc())

	if param.Keyword != "" {
		q = q.
			Where(queries.Role.Name.Like(tools.BuildLikeContains(param.Keyword)))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.Role, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *RoleRepository) FindList(ctx context.Context, idList []int64) ([]*domain.Role, error) {
	if len(idList) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.Role](r.db).
		Where(queries.Role.ID.In(idList...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	list := make([]*domain.Role, 0, len(data))
	for i := range data {
		list = append(list, data[i].ToEntity())
	}

	return list, nil
}

func (r *RoleRepository) FindOne(ctx context.Context, id int64) (*domain.Role, error) {
	m, err := queries.Query[schema.Role](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *RoleRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.Role](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *RoleRepository) NameExist(ctx context.Context, name string) (bool, error) {
	n, err := gorm.G[schema.Role](r.db).
		Where(queries.Role.Name.Eq(name)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *RoleRepository) NameExistExcludeID(ctx context.Context, name string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.Role](r.db).
		Where(queries.Role.Name.Eq(name)).
		Where(queries.Role.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *RoleRepository) Create(ctx context.Context, e domain.Role) error {
	m := schema.Role{Name: e.Name}

	err := gorm.G[schema.Role](r.db).
		Create(ctx, &m)

	return handleError(err)
}

func (r *RoleRepository) Update(ctx context.Context, e domain.Role) error {
	_, err := gorm.G[schema.Role](r.db).
		Where(queries.Role.ID.Eq(e.ID)).
		Set(queries.Role.Name.Set(e.Name)).
		Update(ctx)

	return handleError(err)
}

func (r *RoleRepository) Delete(ctx context.Context, e domain.Role) error {
	_, err := r.enforcer.DeleteRole(GetPolicyRole(e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.Role](r.db).
		Where(queries.Role.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}

func (r *RoleRepository) GrantPermissions(ctx context.Context, role int64, permissions []int64) error {
	policyRole := GetPolicyRole(role)

	_, err := r.enforcer.DeletePermissionsForUser(policyRole)
	if err != nil {
		return handleError(err)
	}

	ps := lo.Map(permissions, func(p int64, _ int) []string {
		return []string{fmt.Sprintf("%d", p)}
	})

	_, err = r.enforcer.AddPermissionsForUser(policyRole, ps...)

	return handleError(err)
}

func (r *RoleRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	pss, err := r.enforcer.GetPermissionsForUser(GetPolicyRole(id))
	if err != nil {
		return nil, handleError(err)
	}

	ps := make([]int64, 0, len(pss))
	for _, s := range pss {
		if len(s) < 2 {
			continue
		}
		i, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			return nil, handleError(err)
		}
		ps = append(ps, i)
	}

	if len(ps) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.Permission](r.db).
		Where(queries.Permission.ID.In(ps...)).
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
