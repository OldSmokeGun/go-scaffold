package usecase

import (
	"context"
	"time"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/pkg/authtoken"
)

var _ AccountUseCaseInterface = (*AccountUseCase)(nil)

type AccountUseCaseInterface interface {
	Login(ctx context.Context, user domain.User) (string, error)
	Logout(ctx context.Context, user domain.User) error
}

type AccountUseCase struct {
	repo repository.UserRepositoryInterface
}

func NewAccountUseCase(
	repo repository.UserRepositoryInterface,
) *AccountUseCase {
	return &AccountUseCase{
		repo: repo,
	}
}

func (c AccountUseCase) Login(ctx context.Context, user domain.User) (string, error) {
	return authtoken.Generate(user.ID, user.Salt, time.Now().Unix()), nil
}

func (c AccountUseCase) Logout(ctx context.Context, user domain.User) error {
	user.RefreshSalt()

	if _, err := c.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
