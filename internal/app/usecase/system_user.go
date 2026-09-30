package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/pkg/pagenation"
)

var _ SystemUserUseCaseInterface = (*SystemUserUseCase)(nil)

type SystemUserUseCaseInterface interface {
	Create(ctx context.Context, user domain.SystemUser) (*domain.SystemUser, error)
	Update(ctx context.Context, user domain.SystemUser) (*domain.SystemUser, error)
	Delete(ctx context.Context, user domain.SystemUser) error
	Detail(ctx context.Context, id int64) (*domain.SystemUser, error)
	List(ctx context.Context, param SystemUserListParam) ([]*domain.SystemUser, int64, error)
	AssignRoles(ctx context.Context, user int64, roles []int64) error
	GetRoles(ctx context.Context, user int64) ([]*domain.SystemRole, error)
	GetRolesByUsers(ctx context.Context, users []int64) (map[int64][]*domain.SystemRole, error)
	GetPermissions(ctx context.Context, user int64) ([]*domain.SystemPermission, error)
	UpdatePassword(ctx context.Context, user domain.SystemUser) error
}

type SystemUserUseCase struct {
	repo repository.SystemUserRepositoryInterface
}

func NewSystemUserUseCase(
	repo repository.SystemUserRepositoryInterface,
) *SystemUserUseCase {
	return &SystemUserUseCase{
		repo: repo,
	}
}

func (c *SystemUserUseCase) Create(ctx context.Context, user domain.SystemUser) (*domain.SystemUser, error) {
	return c.repo.Create(ctx, user)
}

func (c *SystemUserUseCase) Update(ctx context.Context, user domain.SystemUser) (*domain.SystemUser, error) {
	return c.repo.Update(ctx, user)
}

func (c *SystemUserUseCase) Delete(ctx context.Context, user domain.SystemUser) error {
	return c.repo.Delete(ctx, user)
}

func (c *SystemUserUseCase) Detail(ctx context.Context, id int64) (*domain.SystemUser, error) {
	return c.repo.FindOne(ctx, id)
}

type SystemUserListParam struct {
	Keyword string
	pagenation.Param
}

func (c *SystemUserUseCase) List(ctx context.Context, param SystemUserListParam) ([]*domain.SystemUser, int64, error) {
	return c.repo.Filter(ctx, repository.SystemUserFindListParam{
		Keyword: param.Keyword,
		Param:   param.Param,
	})
}

func (c *SystemUserUseCase) AssignRoles(ctx context.Context, user int64, roles []int64) error {
	return c.repo.AssignRoles(ctx, user, roles)
}

func (c *SystemUserUseCase) GetRoles(ctx context.Context, user int64) ([]*domain.SystemRole, error) {
	return c.repo.GetRoles(ctx, user)
}

func (c *SystemUserUseCase) GetRolesByUsers(ctx context.Context, users []int64) (map[int64][]*domain.SystemRole, error) {
	return c.repo.GetRolesByUsers(ctx, users)
}

func (c *SystemUserUseCase) UpdatePassword(ctx context.Context, user domain.SystemUser) error {
	return c.repo.UpdatePassword(ctx, user)
}

func (c *SystemUserUseCase) GetPermissions(ctx context.Context, user int64) ([]*domain.SystemPermission, error) {
	return c.repo.GetPermissions(ctx, user)
}
