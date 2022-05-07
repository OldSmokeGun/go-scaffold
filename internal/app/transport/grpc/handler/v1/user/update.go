package user

import (
	"context"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/service/user"
)

func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	svcReq := user.UpdateRequest{
		Id:    req.Id,
		Name:  req.Name,
		Age:   int8(req.Age),
		Phone: req.Phone,
	}

	_, err := h.service.Update(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}
