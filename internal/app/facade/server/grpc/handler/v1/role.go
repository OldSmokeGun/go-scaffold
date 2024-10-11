package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type RoleHandler struct {
	v1.UnimplementedRoleServer
	logger         *slog.Logger
	roleController *controller.RoleController
}

func NewRoleHandler(
	logger *slog.Logger,
	roleController *controller.RoleController,
) *RoleHandler {
	return &RoleHandler{
		logger:         logger,
		roleController: roleController,
	}
}

func (h *RoleHandler) List(ctx context.Context, req *v1.RoleListRequest) (*v1.RoleListResponse, error) {
	r := controller.RoleListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.roleController.List(ctx, r)
	if err != nil {
		h.logger.Error("call RoleController.List method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.RoleInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.RoleInfo{
			Id:   item.ID,
			Name: item.Name,
		})
	}

	return &v1.RoleListResponse{Items: items}, nil
}

func (h *RoleHandler) Create(ctx context.Context, req *v1.RoleCreateRequest) (*v1.RoleCreateResponse, error) {
	r := controller.RoleCreateRequest{
		RoleAttr: controller.RoleAttr{
			Name: req.Name,
		},
	}

	if err := h.roleController.Create(ctx, r); err != nil {
		h.logger.Error("call RoleController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.RoleCreateResponse{}, nil
}

func (h *RoleHandler) Update(ctx context.Context, req *v1.RoleUpdateRequest) (*v1.RoleUpdateResponse, error) {
	r := controller.RoleUpdateRequest{
		ID: req.Id,
		RoleAttr: controller.RoleAttr{
			Name: req.Name,
		},
	}

	if err := h.roleController.Update(ctx, r); err != nil {
		h.logger.Error("call RoleController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.RoleUpdateResponse{}, nil
}

func (h *RoleHandler) Detail(ctx context.Context, req *v1.RoleDetailRequest) (*v1.RoleInfo, error) {
	ret, err := h.roleController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call RoleController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.RoleInfo{
		Id:   ret.ID,
		Name: ret.Name,
	}, nil
}

func (h *RoleHandler) Delete(ctx context.Context, req *v1.RoleDeleteRequest) (*v1.RoleDeleteResponse, error) {
	if err := h.roleController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call RoleController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.RoleDeleteResponse{}, nil
}

func (h *RoleHandler) GrantPermissions(ctx context.Context, req *v1.RoleGrantPermissionsRequest) (*v1.RoleGrantPermissionsResponse, error) {
	r := controller.RoleGrantPermissionsRequest{
		Role:        req.Role,
		Permissions: req.Permissions,
	}

	if err := h.roleController.GrantPermissions(ctx, r); err != nil {
		h.logger.Error("call RoleController.GrantPermissions method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.RoleGrantPermissionsResponse{}, nil
}

func (h *RoleHandler) GetPermissions(ctx context.Context, req *v1.RoleGetPermissionsRequest) (*v1.RoleGetPermissionsResponse, error) {
	list, err := h.roleController.GetPermissions(ctx, req.Id)
	if err != nil {
		h.logger.Error("call RoleController.GetPermissions method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.PermissionInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.PermissionInfo{
			Id:       item.ID,
			Key:      item.Key,
			Name:     item.Name,
			Desc:     item.Desc,
			ParentID: item.ParentID,
		})
	}

	return &v1.RoleGetPermissionsResponse{Items: items}, nil
}
