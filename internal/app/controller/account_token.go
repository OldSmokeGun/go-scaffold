package controller

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	"go-scaffold/internal/app/service"
	"go-scaffold/internal/app/usecase"
	berr "go-scaffold/internal/errors"
)

type AccountTokenController struct {
	uc   usecase.AccountUseCaseInterface
	repo repository.UserRepositoryInterface
}

func NewAccountTokenController(
	uc usecase.AccountUseCaseInterface,
	repo repository.UserRepositoryInterface,
) *AccountTokenController {
	return &AccountTokenController{
		uc:   uc,
		repo: repo,
	}
}

func (c *AccountTokenController) ValidateToken(ctx context.Context, token string) (*domain.UserProfile, error) {
	claims, err := service.ParseAccountTokenUnverified(token)
	if err != nil {
		return nil, err
	}

	user, err := c.repo.FindOne(ctx, claims.Data.UserID)
	if repository.IsNotFound(err) {
		return nil, berr.ErrInvalidAuthorized.WithError(err)
	} else if err != nil {
		return nil, err
	}

	tokenService := service.NewAccountTokenService(user.Salt)
	_, err = tokenService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return user.ToProfile(), nil
}

func (c *AccountTokenController) RefreshToken(ctx context.Context, userProfile domain.UserProfile, token string) (string, error) {
	claims, err := service.ParseAccountTokenUnverified(token)
	if err != nil {
		return "", err
	}

	expireDuration := claims.ExpiresAt.Sub(time.Now())
	if expireDuration <= 0 {
		return "", berr.ErrInvalidAuthorized.WithError(errors.WithStack(jwt.ErrTokenExpired))
	}
	if expireDuration > domain.AccountTokenRefreshDuration {
		return token, nil
	}

	user := domain.User{
		ID:       userProfile.ID,
		Username: userProfile.Username,
		Nickname: userProfile.Nickname,
		Phone:    userProfile.Phone,
		Salt:     uuid.New().String(),
	}

	token, err = c.uc.Login(ctx, user)
	if err != nil {
		return "", err
	}

	return token, nil
}
