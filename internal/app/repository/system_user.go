package repository

import (
	"context"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository/queries"
	"go-scaffold/internal/app/repository/schema"
	"go-scaffold/internal/app/repository/tools"
	igorm "go-scaffold/internal/pkg/gorm"
	"go-scaffold/internal/pkg/pagenation"
)

var _ SystemUserRepositoryInterface = (*SystemUserRepository)(nil)

type (
	SystemUserFindListParam struct {
		Keyword string
		pagenation.Param
	}

	SystemUserRepositoryInterface interface {
		Filter(ctx context.Context, param SystemUserFindListParam) ([]*domain.SystemUser, int64, error)
		FindOne(ctx context.Context, id int64) (*domain.SystemUser, error)
		FindOneByUsername(ctx context.Context, username string) (*domain.SystemUser, error)
		Exist(ctx context.Context, id int64) (bool, error)
		UsernameExist(ctx context.Context, username string) (bool, error)
		UsernameExistExcludeID(ctx context.Context, username string, excludeID int64) (bool, error)
		Create(ctx context.Context, e domain.SystemUser) (*domain.SystemUser, error)
		Update(ctx context.Context, e domain.SystemUser) (*domain.SystemUser, error)
		Delete(ctx context.Context, e domain.SystemUser) error
		AssignRoles(ctx context.Context, user int64, roles []int64) error
		GetRoles(ctx context.Context, id int64) ([]*domain.SystemRole, error)
		GetRolesByUsers(ctx context.Context, ids []int64) (map[int64][]*domain.SystemRole, error)
		GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error)
		UpdatePassword(ctx context.Context, e domain.SystemUser) error
	}
)

type SystemUserRepository struct {
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewSystemUserRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *SystemUserRepository {
	return &SystemUserRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *SystemUserRepository) Filter(ctx context.Context, param SystemUserFindListParam) ([]*domain.SystemUser, int64, error) {
	build := func() gorm.ChainInterface[schema.SystemUser] {
		q := gorm.G[schema.SystemUser](r.db).
			Order(queries.SystemUser.UpdatedAt.Desc())

		if param.Keyword != "" {
			kw := tools.BuildLikeContains(param.Keyword)
			q = q.
				Where(queries.SystemUser.Username.Like(kw)).
				Or(queries.SystemUser.Nickname.Like(kw)).
				Or(queries.SystemUser.Phone.Like(kw))
		}

		return q
	}

	total, err := build().Count(ctx, "id")
	if err != nil {
		return nil, 0, handleError(err)
	}

	list, err := pagenation.Apply(build(), param.Param).Find(ctx)
	if err != nil {
		return nil, 0, handleError(err)
	}

	entities := make([]*domain.SystemUser, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, total, nil
}

func (r *SystemUserRepository) FindOne(ctx context.Context, id int64) (*domain.SystemUser, error) {
	m, err := queries.Query[schema.SystemUser](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemUserRepository) FindOneByUsername(ctx context.Context, username string) (*domain.SystemUser, error) {
	m, err := gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.Username.Eq(username)).
		Take(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemUserRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.SystemUser](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *SystemUserRepository) UsernameExist(ctx context.Context, username string) (bool, error) {
	n, err := gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.Username.Eq(username)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemUserRepository) UsernameExistExcludeID(ctx context.Context, username string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.Username.Eq(username)).
		Where(queries.SystemUser.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *SystemUserRepository) Create(ctx context.Context, e domain.SystemUser) (*domain.SystemUser, error) {
	m := schema.SystemUser{
		Username: e.Username,
		Password: string(e.Password),
		Nickname: e.Nickname,
		Phone:    e.Phone,
		Salt:     e.Salt,
	}

	if err := gorm.G[schema.SystemUser](r.db).
		Create(ctx, &m); err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *SystemUserRepository) Update(ctx context.Context, e domain.SystemUser) (*domain.SystemUser, error) {
	_, err := gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.ID.Eq(e.ID)).
		Set(
			queries.SystemUser.Username.Set(e.Username),
			queries.SystemUser.Password.Set(string(e.Password)),
			queries.SystemUser.Nickname.Set(e.Nickname),
			queries.SystemUser.Phone.Set(e.Phone),
			queries.SystemUser.Salt.Set(e.Salt),
		).
		Update(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	m := schema.SystemUser{
		BaseModel: igorm.BaseModel{ID: e.ID},
		Username:  e.Username,
		Password:  string(e.Password),
		Nickname:  e.Nickname,
		Phone:     e.Phone,
		Salt:      e.Salt,
	}

	return m.ToEntity(), nil
}

func (r *SystemUserRepository) Delete(ctx context.Context, e domain.SystemUser) error {
	_, err := r.enforcer.DeleteUser(GetPolicyUser(e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}

func (r *SystemUserRepository) AssignRoles(ctx context.Context, user int64, roles []int64) error {
	policyUser := GetPolicyUser(user)

	_, err := r.enforcer.DeleteRolesForUser(policyUser)
	if err != nil {
		return handleError(err)
	}

	rs := lo.Map(roles, func(roleID int64, _ int) string {
		return GetPolicyRole(roleID)
	})

	_, err = r.enforcer.AddRolesForUser(policyUser, rs)

	return handleError(err)
}

func (r *SystemUserRepository) GetRoles(ctx context.Context, id int64) ([]*domain.SystemRole, error) {
	rss, err := r.enforcer.GetRolesForUser(GetPolicyUser(id))
	if err != nil {
		return nil, handleError(err)
	}

	rs := make([]int64, 0, len(rss))
	for _, s := range rss {
		i, err := FromPolicyRole(s)
		if err != nil {
			return nil, err
		}
		rs = append(rs, i)
	}

	if len(rs) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.ID.In(rs...)).
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

func (r *SystemUserRepository) GetRolesByUsers(ctx context.Context, ids []int64) (map[int64][]*domain.SystemRole, error) {
	result := make(map[int64][]*domain.SystemRole, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	userRoleIDs := make(map[int64][]int64, len(ids))
	roleIDSet := make(map[int64]struct{})
	for _, id := range ids {
		rss, err := r.enforcer.GetRolesForUser(GetPolicyUser(id))
		if err != nil {
			return nil, handleError(err)
		}
		for _, item := range rss {
			roleID, err := FromPolicyRole(item)
			if err != nil {
				return nil, err
			}
			userRoleIDs[id] = append(userRoleIDs[id], roleID)
			roleIDSet[roleID] = struct{}{}
		}
	}
	if len(roleIDSet) == 0 {
		return result, nil
	}

	roleIDs := make([]int64, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}

	data, err := gorm.G[schema.SystemRole](r.db).
		Where(queries.SystemRole.ID.In(roleIDs...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	byID := make(map[int64]*domain.SystemRole, len(data))
	for i := range data {
		byID[data[i].ID] = data[i].ToEntity()
	}
	for userID, assignedRoleIDs := range userRoleIDs {
		roles := make([]*domain.SystemRole, 0, len(assignedRoleIDs))
		for _, roleID := range assignedRoleIDs {
			if role, ok := byID[roleID]; ok {
				roles = append(roles, role)
			}
		}
		result[userID] = roles
	}

	return result, nil
}

func (r *SystemUserRepository) UpdatePassword(ctx context.Context, e domain.SystemUser) error {
	_, err := gorm.G[schema.SystemUser](r.db).
		Where(queries.SystemUser.ID.Eq(e.ID)).
		Set(
			queries.SystemUser.Password.Set(string(e.Password)),
			queries.SystemUser.Salt.Set(e.Salt),
		).
		Update(ctx)

	return handleError(err)
}

func (r *SystemUserRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error) {
	// rbac_model.conf treats user_1 as allowed for every object.
	if GetPolicyUser(id) == "user_1" {
		data, err := gorm.G[schema.SystemPermission](r.db).Find(ctx)
		if err != nil {
			return nil, handleError(err)
		}
		return permissionEntities(data), nil
	}

	pss, err := r.enforcer.GetImplicitPermissionsForUser(GetPolicyUser(id))
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

	ps = lo.Uniq(ps)
	if len(ps) == 0 {
		return nil, nil
	}

	data, err := gorm.G[schema.SystemPermission](r.db).
		Where(queries.SystemPermission.ID.In(ps...)).
		Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return permissionEntities(data), nil
}

func permissionEntities(data []schema.SystemPermission) []*domain.SystemPermission {
	list := make([]*domain.SystemPermission, 0, len(data))
	for i := range data {
		list = append(list, data[i].ToEntity())
	}
	return list
}
