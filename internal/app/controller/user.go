package controller

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	"go-scaffold/pkg/validator"
)

type UserController struct {
	uc       usecase.UserUseCaseInterface
	userRepo repository.UserRepositoryInterface
	roleRepo repository.RoleRepositoryInterface
}

func NewUserController(
	uc usecase.UserUseCaseInterface,
	userRepo repository.UserRepositoryInterface,
	roleRepo repository.RoleRepositoryInterface,
) *UserController {
	return &UserController{
		uc:       uc,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

type UserAttr struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

func (r UserAttr) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
			validation.Length(8, 16).Error("username must be 8 ~ 16 characters"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 18).Error("password must be 8 ~ 18 characters"),
			validation.By(validator.PasswordComplexity),
		),
		validation.Field(&r.Nickname,
			validation.Required.Error("nickname is required"),
			validation.Length(8, 16).Error("nickname must be 8 ~ 16 characters"),
		),
		validation.Field(&r.Phone,
			validation.Required.Error("phone is required"),
			validation.By(validator.IsPhoneNumber),
		),
	)
}

type UserCreateRequest struct {
	UserAttr
}

func (r UserCreateRequest) toEntity() domain.User {
	return domain.User{
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
	}
}

func (c *UserController) Create(ctx context.Context, req UserCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	exist, err := c.userRepo.UsernameExist(ctx, req.Username)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("username already exist").WithError(errors.New("username already exist"))
	}

	_, err = c.uc.Create(ctx, req.toEntity())
	return err
}

type UserUpdateRequest struct {
	ID int64 `json:"id"`
	UserAttr
}

func (r UserUpdateRequest) toEntity() domain.User {
	return domain.User{
		ID:       r.ID,
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
	}
}

func (r UserUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.UserAttr),
	)
}

func (c *UserController) Update(ctx context.Context, req UserUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	_, err := c.userRepo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	exist, err := c.userRepo.UsernameExistExcludeID(ctx, req.Username, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.WithMsg("User name already exist").WithError(errors.New("name already exist"))
	}

	_, err = c.uc.Update(ctx, req.toEntity())
	return err
}

func (c *UserController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	role, err := c.userRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}

	return c.uc.Delete(ctx, *role)
}

func (c *UserController) Detail(ctx context.Context, id int64) (*domain.User, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	user, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

type UserListRequest struct {
	Keyword string
}

func (c *UserController) List(ctx context.Context, req UserListRequest) ([]*domain.User, error) {
	param := usecase.UserListParam{Keyword: req.Keyword}
	return c.uc.List(ctx, param)
}

type UserAssignRoleRequest struct {
	User  int64
	Roles []int64
}

func (r UserAssignRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.User, validation.Required.Error("user is required")),
		validation.Field(&r.Roles, validation.Required.Error("no roles that will be assigned")),
	)
}

func (c *UserController) AssignRoles(ctx context.Context, req UserAssignRoleRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	if err := c.validateRolesExist(ctx, req.Roles); err != nil {
		return err
	}

	return c.uc.AssignRoles(ctx, req.User, req.Roles)
}

func (c *UserController) GetRoles(ctx context.Context, id int64) ([]*domain.Role, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	return c.uc.GetRoles(ctx, id)
}

func (c *UserController) validateRolesExist(ctx context.Context, roles []int64) error {
	list, err := c.roleRepo.FindList(ctx, roles)
	if err != nil {
		return err
	}
	roleList := lo.Map(list, func(item *domain.Role, index int) int64 {
		return item.ID
	})

	diffs, _ := lo.Difference(roles, roleList)
	if len(diffs) > 0 {
		return berr.ErrBadCall.WithMsg(fmt.Sprintf("roles %v not exist", diffs)).WithError(errors.New("role not exist"))
	}
	return nil
}
