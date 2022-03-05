package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/pkg/responsex"
	"gorm.io/gorm"
)

func (s *Service) Detail(ctx context.Context, req *pb.DetailRequest) (*pb.DetailReply, error) {
	u, err := s.repo.FindOneByID(
		context.TODO(),
		req.Id,
		[]string{"*"},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		s.logger.Error(err.Error())
		return nil, ErrDataQueryFailed
	}

	result := new(pb.DetailReply)
	if err = copier.Copy(result, u); err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
