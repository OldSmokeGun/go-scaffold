package repository

import (
	"context"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	ient "go-scaffold/internal/pkg/ent"
	"go-scaffold/internal/pkg/ent/ent"
	"go-scaffold/internal/pkg/ent/ent/permission"
	"go-scaffold/internal/pkg/ent/ent/role"
	"go-scaffold/internal/pkg/ent/ent/user"
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
	client   *ient.DefaultClient
	enforcer *casbin.Enforcer
}

func NewUserRepository(client *ient.DefaultClient, enforcer *casbin.Enforcer) *UserRepository {
	return &UserRepository{
		client:   client,
		enforcer: enforcer,
	}
}

func (r *UserRepository) Filter(ctx context.Context, param UserFindListParam) ([]*domain.User, error) {
	query := r.client.User.Query()

	if param.Keyword != "" {
		query.Where(
			user.Or(
				user.UsernameContains(param.Keyword),
				user.NicknameContains(param.Keyword),
				user.PhoneContains(param.Keyword),
			),
		)
	}

	list, err := query.
		Order(ent.Desc(user.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	entities := make([]*domain.User, 0, len(list))
	for _, i := range list {
		entities = append(entities, (&userModel{i}).toEntity())
	}

	return entities, nil
}

func (r *UserRepository) FindOne(ctx context.Context, id int64) (*domain.User, error) {
	m, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&userModel{m}).toEntity(), nil
}

func (r *UserRepository) FindOneByUsername(ctx context.Context, username string) (*domain.User, error) {
	m, err := r.client.User.Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&userModel{m}).toEntity(), nil
}

func (r *UserRepository) Exist(ctx context.Context, id int64) (bool, error) {
	exist, err := r.client.User.Query().Where(user.IDEQ(id)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *UserRepository) UsernameExist(ctx context.Context, username string) (bool, error) {
	exist, err := r.client.User.Query().Where(user.UsernameEQ(username)).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *UserRepository) UsernameExistExcludeID(ctx context.Context, username string, excludeID int64) (bool, error) {
	exist, err := r.client.User.Query().Where(
		user.UsernameEQ(username),
		user.IDNEQ(excludeID),
	).Exist(ctx)
	return exist, errors.WithStack(handleError(err))
}

func (r *UserRepository) Create(ctx context.Context, e domain.User) (*domain.User, error) {
	m, err := r.client.User.Create().
		SetUsername(e.Username).
		SetPassword(string(e.Password)).
		SetNickname(e.Nickname).
		SetPhone(e.Phone).
		SetSalt(e.Salt).
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&userModel{m}).toEntity(), nil
}

func (r *UserRepository) Update(ctx context.Context, e domain.User) (*domain.User, error) {
	m, err := r.client.User.
		UpdateOneID(e.ID).
		SetUsername(e.Username).
		SetPassword(string(e.Password)).
		SetNickname(e.Nickname).
		SetPhone(e.Phone).
		SetSalt(e.Salt).
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}
	return (&userModel{m}).toEntity(), nil
}

func (r *UserRepository) Delete(ctx context.Context, e domain.User) error {
	_, err := r.enforcer.DeleteUser(GetPolicyUser(e.ID))
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(r.client.User.DeleteOneID(e.ID).Exec(ctx))
}

func (r *UserRepository) AssignRoles(ctx context.Context, user int64, roles []int64) error {
	policyUser := GetPolicyUser(user)

	_, err := r.enforcer.DeleteRolesForUser(policyUser)
	if err != nil {
		return errors.WithStack(handleError(err))
	}

	rs := lo.Map(roles, func(r int64, index int) string {
		return GetPolicyRole(r)
	})

	_, err = r.enforcer.AddRolesForUser(policyUser, rs)
	return errors.WithStack(handleError(err))
}

func (r *UserRepository) GetRoles(ctx context.Context, id int64) ([]*domain.Role, error) {
	rss, err := r.enforcer.GetRolesForUser(GetPolicyUser(id))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs := make([]int64, 0, len(rss))
	for _, s := range rss {
		i, err := FromPolicyRole(s)
		if err != nil {
			return nil, err
		}
		rs = append(rs, i)
	}

	data, err := r.client.Role.Query().
		Where(role.IDIn(rs...)).
		All(ctx)
	if err != nil {
		return nil, errors.WithStack(handleError(err))
	}

	list := make([]*domain.Role, 0, len(data))
	for _, item := range data {
		list = append(list, (&roleModel{item}).toEntity())
	}
	return list, err
}

func (r *UserRepository) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	pss, err := r.enforcer.GetImplicitPermissionsForUser(GetPolicyUser(id))
	if err != nil {
		return nil, errors.WithStack(err)
	}

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
	ps = lo.Uniq(ps)

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

type userModel struct {
	*ent.User
}

func (m *userModel) toEntity() *domain.User {
	return &domain.User{
		ID:       m.ID,
		Username: m.Username,
		Password: domain.Password(m.Password),
		Nickname: m.Nickname,
		Phone:    m.Phone,
		Salt:     m.Salt,
	}
}
