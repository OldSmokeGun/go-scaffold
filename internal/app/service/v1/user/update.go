package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/transport/http/pkg/responsex"
	"gorm.io/gorm"
)

func (s *Service) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	u, err := s.repo.FindOneById(
		ctx,
		req.Id,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		s.logger.Error(err)
		return nil, ErrDataQueryFailed
	}

	if err = copier.Copy(u, req); err != nil {
		s.logger.Error(err)
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	if _, err = s.repo.Save(ctx, u); err != nil {
		s.logger.Error(err)
		return nil, ErrDataStoreFailed
	}

	return nil, nil
}
