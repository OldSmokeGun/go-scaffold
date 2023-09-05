package controller

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	"go-scaffold/pkg/validator"
)

type AccountController struct {
	auc      usecase.AccountUseCaseInterface
	uuc      usecase.UserUseCaseInterface
	userRepo repository.UserRepositoryInterface
}

func NewAccountController(
	auc usecase.AccountUseCaseInterface,
	uuc usecase.UserUseCaseInterface,
	userRepo repository.UserRepositoryInterface,
) *AccountController {
	return &AccountController{
		auc:      auc,
		uuc:      uuc,
		userRepo: userRepo,
	}
}

type AccountRegisterRequest struct {
	UserAttr
}

func (r AccountRegisterRequest) toEntity() domain.User {
	return domain.User{
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
		Salt:     uuid.New().String(),
	}
}

type AccountRegisterResponse struct {
	User  *domain.UserProfile `json:"user"`
	Token string              `json:"token"`
}

func (c *AccountController) Register(ctx context.Context, req AccountRegisterRequest) (*AccountRegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	exist, err := c.userRepo.UsernameExist(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, berr.ErrBadCall.WithMsg("username already exist").WithError(errors.New("username already exist"))
	}

	user, err := c.uuc.Create(ctx, req.toEntity())
	if err != nil {
		return nil, err
	}

	token, err := c.auc.Login(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &AccountRegisterResponse{
		user.ToProfile(),
		token,
	}, nil
}

type AccountLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r AccountLoginRequest) Validate() error {
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
	)
}

type AccountLoginResponse struct {
	User  *domain.UserProfile `json:"user"`
	Token string              `json:"token"`
}

func (c *AccountController) Login(ctx context.Context, req AccountLoginRequest) (*AccountLoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	user, err := c.userRepo.FindOneByUsername(ctx, req.Username)
	if repository.IsNotFound(err) {
		return nil, berr.ErrBadCall.WithMsg("username or password is incorrect").WithError(errors.New("username not exist"))
	} else if err != nil {
		return nil, err
	}

	if !user.Password.Verify(domain.Plaintext(req.Password)) {
		return nil, berr.ErrBadCall.WithMsg("username or password is incorrect").WithError(errors.New("password incorrect"))
	}

	token, err := c.auc.Login(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &AccountLoginResponse{
		user.ToProfile(),
		token,
	}, nil
}

func (c *AccountController) Logout(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	user, err := c.userRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithMsg("user not exist").WithError(err)
	} else if err != nil {
		return err
	}

	return c.auc.Logout(ctx, *user)
}

type AccountUpdateProfileRequest struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

func (r AccountUpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.Nickname,
			validation.Required.Error("nickname is required"),
			validation.Length(8, 16).Error("nickname must be 8 ~ 16 characters"),
		),
	)
}

func (c *AccountController) UpdateProfile(ctx context.Context, req AccountUpdateProfileRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	e, err := c.userRepo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return err
	}
	e.Nickname = req.Nickname

	_, err = c.uuc.Update(ctx, *e)
	return err
}

func (c *AccountController) GetProfile(ctx context.Context, id int64) (*domain.UserProfile, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	user, err := c.uuc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.WithError(err)
	} else if err != nil {
		return nil, err
	}

	return user.ToProfile(), nil
}

func (c *AccountController) GetPermissions(ctx context.Context, id int64) ([]*domain.Permission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	return c.uuc.GetPermissions(ctx, id)
}
