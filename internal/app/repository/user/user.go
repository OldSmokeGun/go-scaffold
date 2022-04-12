package user

//go:generate mockgen -source=user.go -destination=user_mock.go -package=user -mock_names=Interface=MockRepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"go-scaffold/internal/app/model"
	"gorm.io/gorm"
	"time"
)

type RepositoryInterface interface {
	// FindByKeyword 根据关键字查询用户列表
	FindByKeyword(ctx context.Context, columns []string, keyword string, order string) ([]*model.User, error)

	// FindOneById 根据 id 查询用户详情
	FindOneById(ctx context.Context, id uint64, columns []string) (*model.User, error)

	// Create 创建用户
	Create(ctx context.Context, user *model.User) (*model.User, error)

	// Save 更新用户信息
	Save(ctx context.Context, user *model.User) (*model.User, error)

	// Delete 删除用户
	Delete(ctx context.Context, user *model.User) error
}

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{
		db:  db,
		rdb: rdb,
	}
}

var (
	cacheKeyFormat = model.User{}.TableName() + "_%d"
	cacheExpire    = 3600
)

func (r *Repository) FindByKeyword(ctx context.Context, columns []string, keyword string, order string) ([]*model.User, error) {
	var users []*model.User
	query := r.db.Select(columns)

	if keyword != "" {
		query.Where("name LIKE ?", "%"+keyword+"%").
			Or("phone LIKE ?", "%"+keyword+"%")
	}

	err := query.Order(order).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) FindOneById(ctx context.Context, id uint64, columns []string) (*model.User, error) {
	m := new(model.User)

	cacheValue, err := r.rdb.Get(
		ctx,
		fmt.Sprintf(cacheKeyFormat, id),
	).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
	}

	if cacheValue != nil {
		if err = jsoniter.Unmarshal(cacheValue, m); err != nil {
			return nil, err
		}

		return m, nil
	}

	err = r.db.Select(columns).Where("id = ?", id).Take(m).Error
	if err != nil {
		return nil, err
	}

	cacheValue, err = jsoniter.Marshal(m)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, id),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *Repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	cacheValue, err := jsoniter.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.Id),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	cacheValue, err := jsoniter.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = r.rdb.Set(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.Id),
		string(cacheValue),
		time.Duration(cacheExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, user *model.User) error {
	if err := r.db.Delete(user).Error; err != nil {
		return err
	}

	err := r.rdb.Del(
		ctx,
		fmt.Sprintf(cacheKeyFormat, user.Id),
	).Err()
	if err != nil {
		return err
	}

	return nil
}
