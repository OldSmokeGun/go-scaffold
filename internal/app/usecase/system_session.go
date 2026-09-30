package usecase

import (
	"context"
	"time"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/pkg/authtoken"
)

var _ SystemSessionUseCaseInterface = (*SystemSessionUseCase)(nil)

type SystemSessionUseCaseInterface interface {
	Login(ctx context.Context, user domain.SystemUser) (string, error)
	Logout(ctx context.Context, user domain.SystemUser) error
}

type SystemSessionUseCase struct {
	repo repository.SystemUserRepositoryInterface
}

func NewSystemSessionUseCase(
	repo repository.SystemUserRepositoryInterface,
) *SystemSessionUseCase {
	return &SystemSessionUseCase{
		repo: repo,
	}
}

func (c SystemSessionUseCase) Login(ctx context.Context, user domain.SystemUser) (string, error) {
	return authtoken.Generate(user.ID, user.Salt, time.Now().Unix()), nil
}

func (c SystemSessionUseCase) Logout(ctx context.Context, user domain.SystemUser) error {
	user.RefreshSalt()

	if _, err := c.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
