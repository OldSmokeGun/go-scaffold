package user

import (
	"context"
	"errors"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"gorm.io/gorm"
)

func (s *Service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
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

	if err = s.repo.Delete(ctx, u); err != nil {
		s.logger.Error(err)
		return nil, ErrDataDeleteFailed
	}

	return &pb.DeleteResponse{}, nil
}
