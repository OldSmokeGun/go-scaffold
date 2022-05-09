package user

import (
	"context"
	"go-scaffold/internal/app/service/user"
	pb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/user"
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
