package domain

import (
	"context"

	"go-scaffold/pkg/validator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

// User 用户实体
type User struct {
	ID    ID
	Name  string // 名称
	Age   int8   // 年龄
	Phone string // 手机号码
}

func (u User) ValidateWithContext(ctx context.Context) error {
	return errors.WithStack(validation.ValidateStruct(&u,
		validation.Field(&u.ID),
		validation.Field(&u.Name, validation.Required.Error("名称不能为空")),
		validation.Field(&u.Phone, validation.By(validator.IsPhoneNumber)),
	))
}

type (
	// UserListParam 用户列表参数
	UserListParam struct {
		Keyword string
	}

	// UserUseCase 用例
	UserUseCase interface {
		// Create 新增用户
		Create(ctx context.Context, user User) error
		// Update 更新用户
		Update(ctx context.Context, user User) error
		// Delete 删除用户
		Delete(ctx context.Context, id ID) error
		// Detail 用户详情
		Detail(ctx context.Context, id ID) (*User, error)
		// List 用户列表
		List(ctx context.Context, param UserListParam) ([]*User, error)
	}
)

type (
	// FindUserListParam 列表查询参数
	FindUserListParam struct {
		Keyword string
	}

	// UserRepository 仓储
	UserRepository interface {
		// FindList 列表查询
		FindList(ctx context.Context, param FindUserListParam) ([]*User, error)
		// FindOneByID 根据 ID 查询详情
		FindOneByID(ctx context.Context, id ID) (*User, error)
		// Exist 数据是否存在
		Exist(ctx context.Context, id ID) (bool, error)
		// Create 新增数据
		Create(ctx context.Context, user User) error
		// Update 更新数据
		Update(ctx context.Context, user User) error
		// Delete 删除数据
		Delete(ctx context.Context, user User) error
	}
)
