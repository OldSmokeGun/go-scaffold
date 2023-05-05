package repository

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/pkg/ent/ent"
	"go-scaffold/internal/app/pkg/ent/ent/user"
)

var _ domain.UserRepository = (*UserRepository)(nil)

// UserRepository 用户仓储
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository 构造用户仓储
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// FindListParam 列表查询参数
type FindListParam struct {
	Keyword string
}

// FindList 列表查询
func (r *UserRepository) FindList(ctx context.Context, param domain.FindUserListParam) ([]*domain.User, error) {
	query := r.client.User.Query()

	if param.Keyword != "" {
		query.Where(user.NameContains(param.Keyword))
	}

	list, err := query.
		Order(ent.Desc(user.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, convertError(err)
	}

	entities := make([]*domain.User, 0, len(list))
	for _, i := range list {
		entities = append(entities, toUserEntity(*i))
	}

	return entities, nil
}

// FindOneByID 根据 id 查询详情
func (r *UserRepository) FindOneByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	m, err := r.client.User.Get(ctx, id.Int64())
	if err != nil {
		return nil, convertError(err)
	}
	return toUserEntity(*m), nil
}

// Exist 数据是否存在
func (r *UserRepository) Exist(ctx context.Context, id domain.ID) (bool, error) {
	exist, err := r.client.User.Query().Where(user.IDEQ(id.Int64())).Exist(ctx)
	return exist, convertError(err)
}

// Create 创建数据
func (r *UserRepository) Create(ctx context.Context, m domain.User) error {
	_, err := r.client.User.Create().
		SetName(m.Name).
		SetAge(m.Age).
		SetPhone(m.Phone).
		Save(ctx)
	return convertError(err)
}

// Update 保存数据
func (r *UserRepository) Update(ctx context.Context, m domain.User) error {
	_, err := r.client.User.
		UpdateOneID(m.ID.Int64()).
		SetName(m.Name).
		SetAge(m.Age).
		SetPhone(m.Phone).
		Save(ctx)
	return convertError(err)
}

// Delete 删除数据
func (r *UserRepository) Delete(ctx context.Context, user domain.User) error {
	return convertError(r.client.User.DeleteOneID(user.ID.Int64()).Exec(ctx))
}

func toUserEntity(m ent.User) *domain.User {
	return &domain.User{
		ID:    domain.ID(m.ID),
		Name:  m.Name,
		Age:   m.Age,
		Phone: m.Phone,
	}
}
