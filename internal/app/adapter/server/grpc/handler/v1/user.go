package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
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

func (h *UserHandler) List(ctx context.Context, req *v1.UserListRequest) (*v1.UserListResponse, error) {
	r := controller.UserListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.userController.List(ctx, r)
	if err != nil {
		h.logger.Error("call UserController.List method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.UserInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.UserInfo{
			Id:       item.ID,
			Username: item.Username,
			Nickname: item.Nickname,
			Phone:    item.Phone,
		})
	}

	return &v1.UserListResponse{Items: items}, nil
}

func (h *UserHandler) Create(ctx context.Context, req *v1.UserCreateRequest) (*v1.UserCreateResponse, error) {
	r := controller.UserCreateRequest{
		UserAttr: controller.UserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}

	if err := h.userController.Create(ctx, r); err != nil {
		h.logger.Error("call UserController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.UserCreateResponse{}, nil
}

func (h *UserHandler) Update(ctx context.Context, req *v1.UserUpdateRequest) (*v1.UserUpdateResponse, error) {
	r := controller.UserUpdateRequest{
		ID: req.Id,
		UserAttr: controller.UserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}

	if err := h.userController.Update(ctx, r); err != nil {
		h.logger.Error("call UserController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.UserUpdateResponse{}, nil
}

func (h *UserHandler) Detail(ctx context.Context, req *v1.UserDetailRequest) (*v1.UserInfo, error) {
	ret, err := h.userController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call UserController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.UserInfo{
		Id:       ret.ID,
		Username: ret.Username,
		Nickname: ret.Nickname,
		Phone:    ret.Phone,
	}, nil
}

func (h *UserHandler) Delete(ctx context.Context, req *v1.UserDeleteRequest) (*v1.UserDeleteResponse, error) {
	if err := h.userController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call UserController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.UserDeleteResponse{}, nil
}

func (h *UserHandler) AssignRoles(ctx context.Context, req *v1.UserAssignRolesRequest) (*v1.UserAssignRolesResponse, error) {
	r := controller.UserAssignRoleRequest{
		User:  req.User,
		Roles: req.Roles,
	}

	if err := h.userController.AssignRoles(ctx, r); err != nil {
		h.logger.Error("call UserController.AssignRoles method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.UserAssignRolesResponse{}, nil
}

func (h *UserHandler) GetRoles(ctx context.Context, req *v1.UserGetRolesRequest) (*v1.UserGetRolesResponse, error) {
	list, err := h.userController.GetRoles(ctx, req.Id)
	if err != nil {
		h.logger.Error("call UserController.GetRoles method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.RoleInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.RoleInfo{
			Id:   item.ID,
			Name: item.Name,
		})
	}

	return &v1.UserGetRolesResponse{Items: items}, nil
}
