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
)

var _ UserRepositoryInterface = (*UserRepository)(nil)

type (
	UserFindListParam struct {
		Keyword string
	}

	UserRepositoryInterface interface {
		Filter(ctx context.Context, param UserFindListParam) ([]*domain.User, error)
		FindOne(ctx context.Context, id int64) (*domain.User, error)
		FindOneByUsername(ctx context.Context, username string) (*domain.User, error)
		Exist(ctx context.Context, id int64) (bool, error)
		UsernameExist(ctx context.Context, username string) (bool, error)
		UsernameExistExcludeID(ctx context.Context, username string, excludeID int64) (bool, error)
		Create(ctx context.Context, e domain.User) (*domain.User, error)
		Update(ctx context.Context, e domain.User) (*domain.User, error)
		Delete(ctx context.Context, e domain.User) error
		AssignRoles(ctx context.Context, user int64, roles []int64) error
		GetRoles(ctx context.Context, id int64) ([]*domain.Role, error)
		GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error)
	}
)

type UserRepository struct {
	db       *igorm.DefaultDB
	enforcer *casbin.Enforcer
}

func NewUserRepository(db *igorm.DefaultDB, enforcer *casbin.Enforcer) *UserRepository {
	return &UserRepository{
		db:       db,
		enforcer: enforcer,
	}
}

func (r *UserRepository) Filter(ctx context.Context, param UserFindListParam) ([]*domain.User, error) {
	q := gorm.G[schema.User](r.db).
		Order(queries.User.UpdatedAt.Desc())

	if param.Keyword != "" {
		kw := tools.BuildLikeContains(param.Keyword)
		q = q.
			Where(queries.User.Username.Like(kw)).
			Or(queries.User.Nickname.Like(kw)).
			Or(queries.User.Phone.Like(kw))
	}

	list, err := q.Find(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	entities := make([]*domain.User, 0, len(list))
	for i := range list {
		entities = append(entities, list[i].ToEntity())
	}

	return entities, nil
}

func (r *UserRepository) FindOne(ctx context.Context, id int64) (*domain.User, error) {
	m, err := queries.Query[schema.User](r.db).FindOne(ctx, id)
	if err == nil && m.ID == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) FindOneByUsername(ctx context.Context, username string) (*domain.User, error) {
	m, err := gorm.G[schema.User](r.db).
		Where(queries.User.Username.Eq(username)).
		Take(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) Exist(ctx context.Context, id int64) (bool, error) {
	n, err := queries.Query[schema.User](r.db).Exist(ctx, id)

	return n > 0, handleError(err)
}

func (r *UserRepository) UsernameExist(ctx context.Context, username string) (bool, error) {
	n, err := gorm.G[schema.User](r.db).
		Where(queries.User.Username.Eq(username)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *UserRepository) UsernameExistExcludeID(ctx context.Context, username string, excludeID int64) (bool, error) {
	n, err := gorm.G[schema.User](r.db).
		Where(queries.User.Username.Eq(username)).
		Where(queries.User.ID.Neq(excludeID)).
		Count(ctx, "*")

	return n > 0, handleError(err)
}

func (r *UserRepository) Create(ctx context.Context, e domain.User) (*domain.User, error) {
	m := schema.User{
		Username: e.Username,
		Password: string(e.Password),
		Nickname: e.Nickname,
		Phone:    e.Phone,
		Salt:     e.Salt,
	}

	if err := gorm.G[schema.User](r.db).
		Create(ctx, &m); err != nil {
		return nil, handleError(err)
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) Update(ctx context.Context, e domain.User) (*domain.User, error) {
	_, err := gorm.G[schema.User](r.db).
		Where(queries.User.ID.Eq(e.ID)).
		Set(
			queries.User.Username.Set(e.Username),
			queries.User.Password.Set(string(e.Password)),
			queries.User.Nickname.Set(e.Nickname),
			queries.User.Phone.Set(e.Phone),
			queries.User.Salt.Set(e.Salt),
		).
		Update(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	m := schema.User{
		ID:       e.ID,
		Username: e.Username,
		Password: string(e.Password),
		Nickname: e.Nickname,
		Phone:    e.Phone,
		Salt:     e.Salt,
	}

	return m.ToEntity(), nil
}

func (r *UserRepository) Delete(ctx context.Context, e domain.User) error {
	_, err := r.enforcer.DeleteUser(GetPolicyUser(e.ID))
	if err != nil {
		return handleError(err)
	}

	_, err = gorm.G[schema.User](r.db).
		Where(queries.User.ID.Eq(e.ID)).
		Delete(ctx)

	return handleError(err)
}

func (r *UserRepository) AssignRoles(ctx context.Context, user int64, roles []int64) error {
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

func (r *UserRepository) GetRoles(ctx context.Context, id int64) ([]*domain.Role, error) {
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

	data, err := gorm.G[schema.Role](r.db).
		Where(queries.Role.ID.In(rs...)).
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

func (r *UserRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
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
