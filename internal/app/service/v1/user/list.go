package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/transport/http/pkg/responsex"
)

func (s *Service) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
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

	result := &pb.ListResponse{Items: []*pb.ListItem{}}
	if err = copier.Copy(&result.Items, users); err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New(responsex.ServerErrorCode.String())
	}

	return result, nil
}
