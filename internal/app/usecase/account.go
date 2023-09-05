package usecase

import (
	"context"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/service"
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
	tokenExpire := domain.AccountTokenExpireDuration

	data := service.AccountTokenData{
		UserID: user.ID,
	}
	return service.NewAccountTokenService(user.Salt).Generate(tokenExpire, data)
}

func (c AccountUseCase) Logout(ctx context.Context, user domain.User) error {
	user.RefreshSalt()

	if _, err := c.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
