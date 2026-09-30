package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/pkg/pagenation"
)

type SystemUserHandler struct {
	v1.UnimplementedSystemUserServer
	logger         *slog.Logger
	userController *controller.SystemUserController
}

func NewSystemUserHandler(
	logger *slog.Logger,
	userController *controller.SystemUserController,
) *SystemUserHandler {
	return &SystemUserHandler{
		logger:         logger,
		userController: userController,
	}
}

func (h *SystemUserHandler) List(ctx context.Context, req *v1.SystemUserListRequest) (*v1.SystemUserListResponse, error) {
	r := controller.SystemUserListRequest{
		Keyword: req.Keyword,
		Param:   pagenation.Param{Limit: pagenation.Unlimit},
	}

	list, _, err := h.userController.List(ctx, r)
	if err != nil {
		h.logger.Error("call SystemUserController.List method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.SystemUserInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.SystemUserInfo{
			Id:       item.ID,
			Username: item.Username,
			Nickname: item.Nickname,
			Phone:    item.Phone,
		})
	}

	return &v1.SystemUserListResponse{Items: items}, nil
}

func (h *SystemUserHandler) Create(ctx context.Context, req *v1.SystemUserCreateRequest) (*v1.SystemUserCreateResponse, error) {
	r := controller.SystemUserCreateRequest{
		SystemUserAttr: controller.SystemUserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}

	if err := h.userController.Create(ctx, r); err != nil {
		h.logger.Error("call SystemUserController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemUserCreateResponse{}, nil
}

func (h *SystemUserHandler) Update(ctx context.Context, req *v1.SystemUserUpdateRequest) (*v1.SystemUserUpdateResponse, error) {
	r := controller.SystemUserUpdateRequest{
		ID: req.Id,
		SystemUserAttr: controller.SystemUserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}

	if err := h.userController.Update(ctx, r); err != nil {
		h.logger.Error("call SystemUserController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemUserUpdateResponse{}, nil
}

func (h *SystemUserHandler) Detail(ctx context.Context, req *v1.SystemUserDetailRequest) (*v1.SystemUserInfo, error) {
	ret, err := h.userController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call SystemUserController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemUserInfo{
		Id:       ret.ID,
		Username: ret.Username,
		Nickname: ret.Nickname,
		Phone:    ret.Phone,
	}, nil
}

func (h *SystemUserHandler) Delete(ctx context.Context, req *v1.SystemUserDeleteRequest) (*v1.SystemUserDeleteResponse, error) {
	if err := h.userController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call SystemUserController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemUserDeleteResponse{}, nil
}

func (h *SystemUserHandler) AssignRoles(ctx context.Context, req *v1.SystemUserAssignRolesRequest) (*v1.SystemUserAssignRolesResponse, error) {
	r := controller.SystemUserAssignRoleRequest{
		User:  req.User,
		Roles: req.Roles,
	}

	if err := h.userController.AssignRoles(ctx, r); err != nil {
		h.logger.Error("call SystemUserController.AssignRoles method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemUserAssignRolesResponse{}, nil
}

func (h *SystemUserHandler) GetRoles(ctx context.Context, req *v1.SystemUserGetRolesRequest) (*v1.SystemUserGetRolesResponse, error) {
	list, err := h.userController.GetRoles(ctx, req.Id)
	if err != nil {
		h.logger.Error("call SystemUserController.GetRoles method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.SystemRoleInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.SystemRoleInfo{
			Id:   item.ID,
			Name: item.Name,
		})
	}

	return &v1.SystemUserGetRolesResponse{Items: items}, nil
}
