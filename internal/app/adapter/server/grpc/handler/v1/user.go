package v1

import (
	"context"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/app/domain"

	"golang.org/x/exp/slog"
)

type UserHandler struct {
	v1.UnimplementedUserServer
	logger         *slog.Logger
	userController *controller.UserController
}

func NewUserHandler(
	logger *slog.Logger,
	userController *controller.UserController,
) *UserHandler {
	return &UserHandler{
		logger:         logger,
		userController: userController,
	}
}

// List 用户列表
func (h *UserHandler) List(ctx context.Context, req *v1.UserListRequest) (*v1.UserListResponse, error) {
	param := domain.UserListParam{
		Keyword: req.Keyword,
	}

	list, err := h.userController.List(ctx, param)
	if err != nil {
		h.logger.Error("call controller.UserController.List method error", err)
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.ListItem, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.ListItem{
			Id:    item.ID.Int64(),
			Name:  item.Name,
			Age:   int32(item.Age),
			Phone: item.Phone,
		})
	}

	resp := &v1.UserListResponse{Items: items}

	return resp, nil
}

// Create 创建用户
func (h *UserHandler) Create(ctx context.Context, req *v1.UserCreateRequest) (*v1.UserCreateResponse, error) {
	p := domain.User{
		Name:  req.Name,
		Age:   int8(req.Age),
		Phone: req.Phone,
	}

	if err := h.userController.Create(ctx, p); err != nil {
		h.logger.Error("call controller.UserController.Create method error", err)
		return nil, errors.Wrap(err)
	}

	return &v1.UserCreateResponse{}, nil
}

// Update 更新用户
func (h *UserHandler) Update(ctx context.Context, req *v1.UserUpdateRequest) (*v1.UserUpdateResponse, error) {
	p := domain.User{
		ID:    domain.ID(req.Id),
		Name:  req.Name,
		Age:   int8(req.Age),
		Phone: req.Phone,
	}

	if err := h.userController.Update(ctx, p); err != nil {
		h.logger.Error("call controller.UserController.Update method error", err)
		return nil, errors.Wrap(err)
	}

	return &v1.UserUpdateResponse{}, nil
}

// Detail 用户详情
func (h *UserHandler) Detail(ctx context.Context, req *v1.UserDetailRequest) (*v1.UserDetailResponse, error) {
	ret, err := h.userController.Detail(ctx, domain.ID(req.Id))
	if err != nil {
		h.logger.Error("call controller.UserController.Detail method error", err)
		return nil, errors.Wrap(err)
	}

	resp := &v1.UserDetailResponse{
		Id:    ret.ID.Int64(),
		Name:  ret.Name,
		Age:   int32(ret.Age),
		Phone: ret.Phone,
	}

	return resp, nil
}

// Delete 删除用户
func (h *UserHandler) Delete(ctx context.Context, req *v1.UserDeleteRequest) (*v1.UserDeleteResponse, error) {
	if err := h.userController.Delete(ctx, domain.ID(req.Id)); err != nil {
		h.logger.Error("call controller.UserController.Delete method error", err)
		return nil, errors.Wrap(err)
	}

	return &v1.UserDeleteResponse{}, nil
}
