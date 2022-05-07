package user

import (
	"context"
	pb "go-scaffold/internal/app/api/scaffold/v1/user"
	"go-scaffold/internal/app/service/user"
)

func (h *Handler) Detail(ctx context.Context, req *pb.DetailRequest) (*pb.DetailResponse, error) {
	svcReq := user.DetailRequest{
		Id: req.Id,
	}

	ret, err := h.service.Detail(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	resp := &pb.DetailResponse{
		Id:    ret.Id,
		Name:  ret.Name,
		Age:   int32(ret.Age),
		Phone: ret.Phone,
	}

	return resp, nil
}
