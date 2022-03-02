package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/v1/user"
	"go-scaffold/internal/app/pkg/responsex"
)

func (s *Service) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	users, err := s.repo.FindByKeyword(
		context.TODO(),
		[]string{"*"},
		req.Keyword,
		"updated_at DESC",
	)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrDataQueryFailed
	}

	result := new(pb.ListReply)
	if err = copier.Copy(&result.Users, users); err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
