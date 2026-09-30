package controller

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
	"go-scaffold/pkg/validator"
)

type SystemSessionController struct {
	systemSessionUseCase usecase.SystemSessionUseCaseInterface
	uuc                  usecase.SystemUserUseCaseInterface
	userRepo             repository.SystemUserRepositoryInterface
}

func NewSystemSessionController(
	systemSessionUseCase usecase.SystemSessionUseCaseInterface,
	uuc usecase.SystemUserUseCaseInterface,
	userRepo repository.SystemUserRepositoryInterface,
) *SystemSessionController {
	return &SystemSessionController{
		systemSessionUseCase: systemSessionUseCase,
		uuc:                  uuc,
		userRepo:             userRepo,
	}
}

type SystemSessionRegisterRequest struct {
	SystemUserAttr
}

func (r SystemSessionRegisterRequest) toEntity() domain.SystemUser {
	return domain.SystemUser{
		Username: r.Username,
		Password: domain.Plaintext(r.Password).Encrypt(),
		Nickname: r.Nickname,
		Phone:    r.Phone,
		Salt:     uuid.New().String(),
	}
}

type SystemSessionRegisterResponse struct {
	User  *domain.SystemUserProfile `json:"user"`
	Token string                    `json:"token"`
}

func (c *SystemSessionController) Register(ctx context.Context, req SystemSessionRegisterRequest) (*SystemSessionRegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	exist, err := c.userRepo.UsernameExist(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, berr.ErrResourceConflict.Errorf("username already exist")
	}

	user, err := c.uuc.Create(ctx, req.toEntity())
	if err != nil {
		return nil, err
	}

	token, err := c.systemSessionUseCase.Login(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &SystemSessionRegisterResponse{
		user.ToProfile(),
		token,
	}, nil
}

type SystemSessionLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r SystemSessionLoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 18).Error("password must be 8 ~ 18 characters"),
			validation.By(validator.PasswordComplexity),
		),
	)
}

type SystemSessionLoginResponse struct {
	User  *domain.SystemUserProfile `json:"user"`
	Token string                    `json:"token"`
}

func (c *SystemSessionController) Login(ctx context.Context, req SystemSessionLoginRequest) (*SystemSessionLoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	user, err := c.userRepo.FindOneByUsername(ctx, req.Username)
	if repository.IsNotFound(err) {
		return nil, berr.ErrBadCall.WithMsg("username or password is incorrect").Wrap(err)
	} else if err != nil {
		return nil, err
	}

	if !user.Password.Verify(domain.Plaintext(req.Password)) {
		return nil, berr.ErrBadCall.WithMsg("username or password is incorrect").Wrap(errors.New("password incorrect"))
	}

	token, err := c.systemSessionUseCase.Login(ctx, *user)
	if err != nil {
		return nil, err
	}

	profile := user.ToProfile()
	roles, err := c.uuc.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	profile.Roles = roles

	return &SystemSessionLoginResponse{
		profile,
		token,
	}, nil
}

func (c *SystemSessionController) Logout(ctx context.Context, id int64) error {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	user, err := c.userRepo.FindOne(ctx, id)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.WithMsg("user not exist").Wrap(err)
	} else if err != nil {
		return err
	}

	return c.systemSessionUseCase.Logout(ctx, *user)
}

type SystemSessionUpdateProfileRequest struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

func (r SystemSessionUpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.Nickname,
			validation.Required.Error("nickname is required"),
			validation.Length(8, 16).Error("nickname must be 8 ~ 16 characters"),
		),
	)
}

func (c *SystemSessionController) UpdateProfile(ctx context.Context, req SystemSessionUpdateProfileRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	e, err := c.userRepo.FindOne(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}
	e.Nickname = req.Nickname

	_, err = c.uuc.Update(ctx, *e)
	return err
}

type SystemSessionUpdatePasswordRequest struct {
	ID          int64
	OldPassword string
	Password    string
}

func (r SystemSessionUpdatePasswordRequest) Validate() error {
	if err := validation.Validate(r.ID, validation.Required.Error("id is required")); err != nil {
		return err
	}
	if err := validation.Validate(r.OldPassword, validation.Required.Error("old password is required")); err != nil {
		return err
	}

	return validatePassword(r.Password)
}

func (c *SystemSessionController) UpdatePassword(ctx context.Context, req SystemSessionUpdatePasswordRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.Wrap(err)
	}

	user, err := c.uuc.Detail(ctx, req.ID)
	if repository.IsNotFound(err) {
		return berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return err
	}

	if !user.Password.Verify(domain.Plaintext(req.OldPassword)) {
		return berr.ErrBadCall.WithMsg("old password is incorrect").Wrap(errors.New("old password incorrect"))
	}

	user.Password = domain.Plaintext(req.Password).Encrypt()
	user.RefreshSalt()

	return c.uuc.UpdatePassword(ctx, *user)
}

func (c *SystemSessionController) GetProfile(ctx context.Context, id int64) (*domain.SystemUserProfile, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	user, err := c.uuc.Detail(ctx, id)
	if repository.IsNotFound(err) {
		return nil, berr.ErrResourceNotFound.Wrap(err)
	} else if err != nil {
		return nil, err
	}

	profile := user.ToProfile()
	roles, err := c.uuc.GetRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	profile.Roles = roles

	return profile, nil
}

func (c *SystemSessionController) GetPermissions(ctx context.Context, id int64) ([]*domain.SystemPermission, error) {
	if err := validation.Validate(id, validation.Required.Error("id is required")); err != nil {
		return nil, berr.ErrValidateError.Wrap(err)
	}

	return c.uuc.GetPermissions(ctx, id)
}
