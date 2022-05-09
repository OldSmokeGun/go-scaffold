package user

import (
	"context"
	"go-scaffold/internal/app/service/user"
	pb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/user"
)

func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	svcReq := user.CreateRequest{
		Name:  req.Name,
		Age:   int8(req.Age),
		Phone: req.Phone,
	}

	_, err := h.service.Create(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}
