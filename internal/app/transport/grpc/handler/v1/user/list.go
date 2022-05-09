package user

import (
	"context"
	"go-scaffold/internal/app/service/user"
	pb "go-scaffold/internal/app/transport/grpc/api/scaffold/v1/user"
)

func (h *Handler) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	svcReq := user.ListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.service.List(ctx, svcReq)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.ListItem, 0, len(list))

	for _, item := range list {
		items = append(items, &pb.ListItem{
			Id:    item.Id,
			Name:  item.Name,
			Age:   int32(item.Age),
			Phone: item.Phone,
		})
	}

	resp := &pb.ListResponse{Items: items}

	return resp, nil
}
