package user

import (
	"github.com/go-redis/redis/v8"
	"go-scaffold/internal/app/global"
	"go-scaffold/internal/app/model"
	"gorm.io/gorm"
)

type Repository interface {
	QueryOneByID(id int) *model.User
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
	// TODO
	
	return new(model.User)
}
