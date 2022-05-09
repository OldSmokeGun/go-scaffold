package greet

import (
	"context"
	"go-scaffold/internal/app/service/greet"
	pb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/greet"
)

func (h *Handler) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	svcReq := greet.HelloRequest{Name: req.Name}

	ret, err := h.service.Hello(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	resp := &pb.HelloResponse{
		Msg: ret.Msg,
	}

	return resp, nil
}
