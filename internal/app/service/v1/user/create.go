package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/model"
	"go-scaffold/internal/app/pkg/responsex"
)

func (s *Service) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	m := new(model.User)
	if err := copier.Copy(m, req); err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	if _, err := s.repo.Create(context.TODO(), m); err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDataStoreFailed
	}

	return nil, nil
}
