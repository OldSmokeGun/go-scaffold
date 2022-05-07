package user

import (
	"context"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/service/user"
)

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	svcReq := user.DeleteRequest{
		Id: req.Id,
	}

	if err := h.service.Delete(ctx, svcReq); err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}
