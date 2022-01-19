package user

import (
	"github.com/go-redis/redis/v8"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/model"
	"gorm.io/gorm"
)

type Repository interface {
	// 示例方法，根据业务编写

	FindOneByID(id uint, columns []string) *model.User
	FindOneWhere(where map[string]interface{}, columns []string, order string) *model.User
	FindAll(columns []string, order string) []*model.User
	FindWhere(where map[string]interface{}, columns []string, order string) []*model.User
	Create(user *model.User) (*model.User, error)
	Save(user *model.User) (*model.User, error)
	Delete(user *model.User) error
	Paginate(limit int, columns []string, order string)
}

type repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func New() *repository {
	return &repository{
		db:  global.DB(),
		rdb: global.RedisClient(),
	}
}

func (r repository) QueryOneByID(id int) *model.User {
	panic("implement me")
}

func (r repository) FindOneByID(id uint, columns []string) *model.User {
	panic("implement me")
}

func (r repository) FindOneWhere(where map[string]interface{}, columns []string, order string) *model.User {
	panic("implement me")
}

func (r repository) FindAll(columns []string, order string) []*model.User {
	panic("implement me")
}

func (r repository) FindWhere(where map[string]interface{}, columns []string, order string) []*model.User {
	panic("implement me")
}

func (r repository) Create(user *model.User) (*model.User, error) {
	panic("implement me")
}

func (r repository) Save(user *model.User) (*model.User, error) {
	panic("implement me")
}

func (r repository) Delete(user *model.User) error {
	panic("implement me")
}

func (r repository) Paginate(limit int, columns []string, order string) {
	panic("implement me")
}
