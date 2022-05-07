package user

import (
	"context"
	"go-scaffold/internal/app/model"
	errorsx "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/app/repository/user"
)

// ListRequest 用户列表请求参数
type ListRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
}

// ListItem 用户列表项
type ListItem = DetailResponse

// ListResponse 用户列表响应数据
type ListResponse []*ListItem

// List 用户列表
func (s *Service) List(ctx context.Context, req ListRequest) (ListResponse, error) {
	list, err := s.repo.FindList(
		ctx,
		user.FindListParam{Keyword: req.Keyword},
		[]string{"*"},
		"updated_at DESC",
	)
	if err != nil {
		s.logger.Errorf("%s: %s", model.ErrDataQueryFailed, err)
		return nil, errorsx.ServerError(errorsx.WithMessage(model.ErrDataQueryFailed.Error()))
	}

	resp := make(ListResponse, 0, len(list))

	for _, item := range list {
		resp = append(resp, &ListItem{
			Id:    item.Id,
			Name:  item.Name,
			Age:   item.Age,
			Phone: item.Phone,
		})
	}

	return resp, nil
}
