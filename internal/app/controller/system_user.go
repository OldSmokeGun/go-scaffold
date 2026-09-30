package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
	"go-scaffold/internal/pkg/pagenation"
	"go-scaffold/pkg/validator"
)

type SystemUserController struct {
	uc       usecase.SystemUserUseCaseInterface
	userRepo repository.SystemUserRepositoryInterface
	roleRepo repository.SystemRoleRepositoryInterface
}

func NewSystemUserController(
	uc usecase.SystemUserUseCaseInterface,
	userRepo repository.SystemUserRepositoryInterface,
	roleRepo repository.SystemRoleRepositoryInterface,
) *SystemUserController {
	return &SystemUserController{
		uc:       uc,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

type SystemUserAttr struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

func (r SystemUserAttr) Validate() error {
	if err := r.ValidateProfile(); err != nil {
		return err
	}

	return validatePassword(r.Password)
}

func (r SystemUserAttr) ValidateProfile() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
			validation.Length(8, 16).Error("username must be 8 ~ 16 characters"),
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

func validatePassword(password string) error {
	return validation.Validate(password,
		validation.Required.Error("password is required"),
		validation.Length(8, 18).Error("password must be 8 ~ 18 characters"),
		validation.By(validator.PasswordComplexity),
	)
}

type SystemUserCreateRequest struct {
	SystemUserAttr
}

func (r SystemUserCreateRequest) toEntity() domain.SystemUser {
	return domain.SystemUser{
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
	}
}

func (c *SystemUserController) Create(ctx context.Context, req SystemUserCreateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	exist, err := c.userRepo.UsernameExist(ctx, req.Username)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("username already exist")
	}

	_, err = c.uc.Create(ctx, req.toEntity())
	return err
}

type SystemUserUpdateRequest struct {
	ID int64 `json:"id"`
	SystemUserAttr
}

func (r SystemUserUpdateRequest) toEntity() domain.SystemUser {
	return domain.SystemUser{
		ID:       r.ID,
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
	}
}

func (r SystemUserUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.SystemUserAttr),
	)
}

func (c *SystemUserController) Update(ctx context.Context, req SystemUserUpdateRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	current, err := c.uc.Detail(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	exist, err := c.userRepo.UsernameExistExcludeID(ctx, req.Username, req.ID)
	if err != nil {
		return err
	}
	if exist {
		return berr.ErrBadCall.Errorf("username already exist")
	}

	entity := req.toEntity()
	entity.Salt = current.Salt

	_, err = c.uc.Update(ctx, entity)
	return err
}

func (c *SystemUserController) Delete(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	role, err := c.userRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	return c.uc.Delete(ctx, *role)
}

func (c *SystemUserController) Detail(ctx context.Context, id int64) (*domain.SystemUser, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	user, err := c.uc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

type SystemUserListRequest struct {
	Keyword string
	pagenation.Param
}

type SystemUserListItem struct {
	*domain.SystemUser
	Roles []*domain.SystemRole
}

func (c *SystemUserController) List(ctx context.Context, req SystemUserListRequest) ([]*SystemUserListItem, int64, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	users, total, err := c.uc.List(ctx, usecase.SystemUserListParam{
		Keyword: req.Keyword,
		Param:   req.Param,
	})
	if err != nil {
		return nil, 0, err
	}

	ids := lo.Map(users, func(item *domain.SystemUser, _ int) int64 {
		return item.ID
	})
	rolesByUser, err := c.uc.GetRolesByUsers(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*SystemUserListItem, 0, len(users))
	for _, user := range users {
		roles := rolesByUser[user.ID]
		if roles == nil {
			roles = []*domain.SystemRole{}
		}
		items = append(items, &SystemUserListItem{
			SystemUser: user,
			Roles:      roles,
		})
	}

	return items, total, nil
}

type SystemUserAssignRoleRequest struct {
	User  int64
	Roles []int64
}

func (r SystemUserAssignRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.User, validation.Required.Error("user is required")),
		validation.Field(&r.Roles, validation.Required.Error("no roles that will be assigned")),
	)
}

func (c *SystemUserController) AssignRoles(ctx context.Context, req SystemUserAssignRoleRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	if err := c.validateRolesExist(ctx, req.Roles); err != nil {
		return err
	}

	return c.uc.AssignRoles(ctx, req.User, req.Roles)
}

func (c *SystemUserController) GetRoles(ctx context.Context, id int64) ([]*domain.SystemRole, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	return c.uc.GetRoles(ctx, id)
}

func (c *SystemUserController) validateRolesExist(ctx context.Context, roles []int64) error {
	list, err := c.roleRepo.FindList(ctx, roles)
	if err != nil {
		return err
	}
	roleList := lo.Map(list, func(item *domain.SystemRole, index int) int64 {
		return item.ID
	})

	diffs, _ := lo.Difference(roles, roleList)
	if len(diffs) > 0 {
		return berr.ErrBadCall.Errorf("roles %v not exist", diffs)
	}
	return nil
}
