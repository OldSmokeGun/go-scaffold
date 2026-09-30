package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type SystemRoleHandler struct {
	v1.UnimplementedSystemRoleServer
	logger         *slog.Logger
	roleController *controller.SystemRoleController
}

func NewSystemRoleHandler(
	logger *slog.Logger,
	roleController *controller.SystemRoleController,
) *SystemRoleHandler {
	return &SystemRoleHandler{
		logger:         logger,
		roleController: roleController,
	}
}

func (h *SystemRoleHandler) List(ctx context.Context, req *v1.SystemRoleListRequest) (*v1.SystemRoleListResponse, error) {
	r := controller.SystemRoleListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.roleController.List(ctx, r)
	if err != nil {
		h.logger.Error("call SystemRoleController.List method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.SystemRoleInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.SystemRoleInfo{
			Id:   item.ID,
			Name: item.Name,
		})
	}

	return &v1.SystemRoleListResponse{Items: items}, nil
}

func (h *SystemRoleHandler) Create(ctx context.Context, req *v1.SystemRoleCreateRequest) (*v1.SystemRoleCreateResponse, error) {
	r := controller.SystemRoleCreateRequest{
		SystemRoleAttr: controller.SystemRoleAttr{
			Name: req.Name,
		},
	}

	if err := h.roleController.Create(ctx, r); err != nil {
		h.logger.Error("call SystemRoleController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemRoleCreateResponse{}, nil
}

func (h *SystemRoleHandler) Update(ctx context.Context, req *v1.SystemRoleUpdateRequest) (*v1.SystemRoleUpdateResponse, error) {
	r := controller.SystemRoleUpdateRequest{
		ID: req.Id,
		SystemRoleAttr: controller.SystemRoleAttr{
			Name: req.Name,
		},
	}

	if err := h.roleController.Update(ctx, r); err != nil {
		h.logger.Error("call SystemRoleController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemRoleUpdateResponse{}, nil
}

func (h *SystemRoleHandler) Detail(ctx context.Context, req *v1.SystemRoleDetailRequest) (*v1.SystemRoleInfo, error) {
	ret, err := h.roleController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call SystemRoleController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemRoleInfo{
		Id:   ret.ID,
		Name: ret.Name,
	}, nil
}

func (h *SystemRoleHandler) Delete(ctx context.Context, req *v1.SystemRoleDeleteRequest) (*v1.SystemRoleDeleteResponse, error) {
	if err := h.roleController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call SystemRoleController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemRoleDeleteResponse{}, nil
}

func (h *SystemRoleHandler) GrantPermissions(ctx context.Context, req *v1.SystemRoleGrantPermissionsRequest) (*v1.SystemRoleGrantPermissionsResponse, error) {
	r := controller.SystemRoleGrantPermissionsRequest{
		Role:        req.Role,
		Permissions: req.Permissions,
	}

	if err := h.roleController.GrantPermissions(ctx, r); err != nil {
		h.logger.Error("call SystemRoleController.GrantPermissions method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemRoleGrantPermissionsResponse{}, nil
}

func (h *SystemRoleHandler) GetPermissions(ctx context.Context, req *v1.SystemRoleGetPermissionsRequest) (*v1.SystemRoleGetPermissionsResponse, error) {
	list, err := h.roleController.GetPermissions(ctx, req.Id)
	if err != nil {
		h.logger.Error("call SystemRoleController.GetPermissions method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.SystemPermissionInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.SystemPermissionInfo{
			Id:       item.ID,
			Key:      item.Key,
			Name:     item.Name,
			Desc:     item.Desc,
			ParentID: item.ParentID,
		})
	}

	return &v1.SystemRoleGetPermissionsResponse{Items: items}, nil
}
