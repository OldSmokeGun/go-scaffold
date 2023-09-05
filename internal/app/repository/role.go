package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/app/pkg/ent/ent/permission"
	"go-scaffold/internal/app/pkg/ent/ent/role"
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
	client   *ent.Client
	enforcer *casbin.Enforcer
}

func NewRoleRepository(client *ent.Client, enforcer *casbin.Enforcer) *RoleRepository {
	return &RoleRepository{
		client:   client,
		enforcer: enforcer,
	}
}

func (r *RoleRepository) Filter(ctx context.Context, param RoleFindListParam) ([]*domain.Role, error) {
	query := r.client.Role.Query()

	if param.Keyword != "" {
		query.Where(role.NameContains(param.Keyword))
	}

	list, err := query.
		Order(ent.Desc(role.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	entities := make([]*domain.Role, 0, len(list))
	for _, i := range list {
		entities = append(entities, (&roleModel{i}).toEntity())
	}

	return entities, nil
}

func (r *RoleRepository) FindList(ctx context.Context, idList []int64) ([]*domain.Role, error) {
	data, err := r.client.Role.Query().
		Where(role.IDIn(idList...)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	list := make([]*domain.Role, 0, len(data))
	for _, item := range data {
		list = append(list, (&roleModel{item}).toEntity())
	}
	return list, nil
}

func (r *RoleRepository) FindOne(ctx context.Context, id int64) (*domain.Role, error) {
	m, err := r.client.Role.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&roleModel{m}).toEntity(), nil
}

func (r *RoleRepository) Exist(ctx context.Context, id int64) (bool, error) {
	exist, err := r.client.Role.Query().Where(role.IDEQ(id)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *RoleRepository) NameExist(ctx context.Context, name string) (bool, error) {
	exist, err := r.client.Role.Query().Where(role.NameEQ(name)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *RoleRepository) NameExistExcludeID(ctx context.Context, name string, excludeID int64) (bool, error) {
	exist, err := r.client.Role.Query().Where(
		role.NameEQ(name),
		role.IDNEQ(excludeID),
	).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *RoleRepository) Create(ctx context.Context, e domain.Role) error {
	_, err := r.client.Role.Create().
		SetName(e.Name).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *RoleRepository) Update(ctx context.Context, e domain.Role) error {
	_, err := r.client.Role.
		UpdateOneID(e.ID).
		SetName(e.Name).
		Save(ctx)
	return errors.WithStack(handleError(err))
}

func (r *RoleRepository) Delete(ctx context.Context, e domain.Role) error {
	_, err := r.enforcer.DeleteRole(GetPolicyRole(e.ID))
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(r.client.Role.DeleteOneID(e.ID).Exec(ctx))
}

func (r *RoleRepository) GrantPermissions(ctx context.Context, role int64, permissions []int64) error {
	policyRole := GetPolicyRole(role)

	_, err := r.enforcer.DeletePermissionsForUser(policyRole)
	if err != nil {
		return errors.WithStack(handleError(err))
	}

	ps := lo.Map(permissions, func(p int64, index int) []string {
		return []string{fmt.Sprintf("%d", p)}
	})

	_, err = r.enforcer.AddPermissionsForUser(policyRole, ps...)
	return errors.WithStack(handleError(err))
}

func (r *RoleRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	pss := r.enforcer.GetPermissionsForUser(GetPolicyRole(id))

	ps := make([]int64, 0, len(pss))
	for _, s := range pss {
		if len(s) < 2 {
			continue
		}
		i, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ps = append(ps, i)
	}

	data, err := r.client.Permission.Query().
		Where(permission.IDIn(ps...)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	list := make([]*domain.Permission, 0, len(data))
	for _, item := range data {
		list = append(list, (&permissionModel{item}).toEntity())
	}
	return list, err
}

type roleModel struct {
	*ent.Role
}

func (m *roleModel) toEntity() *domain.Role {
	return &domain.Role{
		ID:   m.ID,
		Name: m.Name,
	}
}
