package controller

import (
	"context"
	"time"

	"go-scaffold/internal/app/domain"
	"go-scaffold/internal/app/repository"
	berr "go-scaffold/internal/errors"
	"go-scaffold/pkg/authtoken"
)

type SystemSessionTokenController struct {
	repo repository.SystemUserRepositoryInterface
}

func NewSystemSessionTokenController(repo repository.SystemUserRepositoryInterface) *SystemSessionTokenController {
	return &SystemSessionTokenController{repo: repo}
}

func (c *SystemSessionTokenController) ValidateToken(ctx context.Context, token string) (*domain.SystemUserProfile, error) {
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
	if !authtoken.InExpireWindow(ts, time.Now(), domain.SystemSessionTokenExpireDuration) {
		return nil, berr.ErrInvalidAuthorized.WithMsg("token is expired")
	}

	return user.ToProfile(), nil
}

// RefreshToken 有效窗口内续期，返回新 token（不更换 salt）。
func (c *SystemSessionTokenController) RefreshToken(ctx context.Context, userProfile domain.SystemUserProfile, token string) (string, error) {
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
