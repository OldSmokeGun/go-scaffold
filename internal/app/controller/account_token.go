package controller

import (
	"context"
	"time"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	berr "go-scaffold/internal/errors"
	"go-scaffold/pkg/authtoken"
)

type AccountTokenController struct {
	repo repository.UserRepositoryInterface
}

func NewAccountTokenController(repo repository.UserRepositoryInterface) *AccountTokenController {
	return &AccountTokenController{repo: repo}
}

func (c *AccountTokenController) ValidateToken(ctx context.Context, token string) (*domain.UserProfile, error) {
	userID, ts, sign, err := authtoken.Parse(token)
	if err != nil {
		return nil, berr.ErrInvalidAuthorized.WithError(err)
	}

	user, err := c.repo.FindOne(ctx, userID)
	if repository.IsNotFound(err) {
		return nil, berr.ErrInvalidAuthorized.WithError(err)
	} else if err != nil {
		return nil, err
	}
	if !authtoken.VerifySign(userID, ts, user.Salt, sign) {
		return nil, berr.ErrInvalidAuthorized.WithMsg("token is invalid")
	}
	if !authtoken.InExpireWindow(ts, time.Now(), domain.AccountTokenExpireDuration) {
		return nil, berr.ErrInvalidAuthorized.WithMsg("token is expired")
	}

	return user.ToProfile(), nil
}

// RefreshToken 有效窗口内续期，返回新 token（不更换 salt）。
func (c *AccountTokenController) RefreshToken(ctx context.Context, userProfile domain.UserProfile, token string) (string, error) {
	_ = token
	if userProfile.ID <= 0 {
		return "", berr.ErrInvalidAuthorized.WithMsg("token is invalid")
	}
	user, err := c.repo.FindOne(ctx, userProfile.ID)
	if repository.IsNotFound(err) {
		return "", berr.ErrInvalidAuthorized.WithError(err)
	}
	if err != nil {
		return "", err
	}

	return authtoken.Generate(user.ID, user.Salt, time.Now().Unix()), nil
}
