package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/transport/http/pkg/responsex"
	"gorm.io/gorm"
)

func (s *Service) Detail(ctx context.Context, req *pb.DetailRequest) (*pb.DetailResponse, error) {
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

	result := new(pb.DetailResponse)
	if err = copier.Copy(result, u); err != nil {
		s.logger.Error(err)
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
