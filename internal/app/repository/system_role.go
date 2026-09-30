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

var _ SystemRoleRepositoryInterface = (*SystemRoleRepository)(nil)

type (
	SystemRoleFindListParam struct {
		Keyword string
	}

	SystemRoleRepositoryInterface interface {
		Filter(ctx context.Context, param SystemRoleFindListParam) ([]*domain.SystemRole, error)
		FindList(ctx context.Context, idList []int64) ([]*domain.SystemRole, error)
		FindOne(ctx context.Context, id int64) (*domain.SystemRole, error)
		Exist(ctx context.Context, id int64) (bool, error)
		NameExist(ctx context.Context, name string) (bool, error)
		NameExistExcludeID(ctx context.Context, name string, excludeID int64) (bool, error)
		Create(ctx context.Context, e domain.SystemRole) error
		Update(ctx context.Context, e domain.SystemRole) error
		Delete(ctx context.Context, e domain.SystemRole) error
		GrantPermissions(ctx context.Context, role int64, permissions []int64) error
		GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error)
	}
)

type SystemRoleRepository struct {
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewSystemRoleRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *SystemRoleRepository {
	return &SystemRoleRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *SystemRoleRepository) Filter(ctx context.Context, param SystemRoleFindListParam) ([]*domain.SystemRole, error) {
	q := gorm.G[schema.SystemRole](r.db).
		Order(queries.SystemRole.UpdatedAt.Desc())

	if param.Keyword != "" {
		q = q.
			Where(queries.SystemRole.Name.Like(tools.BuildLikeContains(param.Keyword)))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.SystemRole, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *SystemRoleRepository) FindList(ctx context.Context, idList []int64) ([]*domain.SystemRole, error) {
	if len(idList) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.ID.In(idList...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	list := make([]*domain.SystemRole, 0, len(data))
	for i := range data {
		list = append(list, data[i].ToEntity())
	}

	return list, nil
}

func (r *SystemRoleRepository) FindOne(ctx context.Context, id int64) (*domain.SystemRole, error) {
	m, err := queries.Query[schema.SystemRole](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemRoleRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.SystemRole](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *SystemRoleRepository) NameExist(ctx context.Context, name string) (bool, error) {
	n, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.Name.Eq(name)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemRoleRepository) NameExistExcludeID(ctx context.Context, name string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.Name.Eq(name)).
		Where(queries.SystemRole.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemRoleRepository) Create(ctx context.Context, e domain.SystemRole) error {
	m := schema.SystemRole{Name: e.Name}

	err := gorm.G[schema.SystemRole](r.db).
		Create(ctx, &m)

	return handleError(err)
}

func (r *SystemRoleRepository) Update(ctx context.Context, e domain.SystemRole) error {
	_, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.ID.Eq(e.ID)).
		Set(queries.SystemRole.Name.Set(e.Name)).
		Update(ctx)

	return handleError(err)
}

func (r *SystemRoleRepository) Delete(ctx context.Context, e domain.SystemRole) error {
	_, err := r.enforcer.DeleteRole(GetPolicyRole(e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}

func (r *SystemRoleRepository) GrantPermissions(ctx context.Context, role int64, permissions []int64) error {
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

func (r *SystemRoleRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error) {
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

	data, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ID.In(ps...)).
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
