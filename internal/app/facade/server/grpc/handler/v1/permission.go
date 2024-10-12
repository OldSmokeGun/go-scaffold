package v1

import (
	"context"
	"log/slog"

	"go-scaffold/internal/app/controller"
	v1 "go-scaffold/internal/app/facade/server/grpc/api/v1"
	"go-scaffold/internal/app/facade/server/grpc/pkg/errors"
)

type PermissionHandler struct {
	v1.UnimplementedPermissionServer
	logger               *slog.Logger
	permissionController *controller.PermissionController
}

func NewPermissionHandler(
	logger *slog.Logger,
	permissionController *controller.PermissionController,
) *PermissionHandler {
	return &PermissionHandler{
		logger:               logger,
		permissionController: permissionController,
	}
}

// List 权限列表
func (h *PermissionHandler) List(ctx context.Context, req *v1.PermissionListRequest) (*v1.PermissionListResponse, error) {
	r := controller.PermissionListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.permissionController.List(ctx, r)
	if err != nil {
		h.logger.Error("call PermissionController.List method error", slog.Any("error", err))
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

	return &v1.PermissionListResponse{Items: items}, nil
}

// Create 权限创建
func (h *PermissionHandler) Create(ctx context.Context, req *v1.PermissionCreateRequest) (*v1.PermissionCreateResponse, error) {
	r := controller.PermissionCreateRequest{
		PermissionAttr: controller.PermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}

	if err := h.permissionController.Create(ctx, r); err != nil {
		h.logger.Error("call PermissionController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.PermissionCreateResponse{}, nil
}

// Update 权限更新
func (h *PermissionHandler) Update(ctx context.Context, req *v1.PermissionUpdateRequest) (*v1.PermissionUpdateResponse, error) {
	r := controller.PermissionUpdateRequest{
		ID: req.Id,
		PermissionAttr: controller.PermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}

	if err := h.permissionController.Update(ctx, r); err != nil {
		h.logger.Error("call PermissionController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.PermissionUpdateResponse{}, nil
}

// Detail 权限详情
func (h *PermissionHandler) Detail(ctx context.Context, req *v1.PermissionDetailRequest) (*v1.PermissionInfo, error) {
	ret, err := h.permissionController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call PermissionController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.PermissionInfo{
		Id:       ret.ID,
		Key:      ret.Key,
		Name:     ret.Name,
		Desc:     ret.Desc,
		ParentID: ret.ParentID,
	}, nil
}

// Delete 权限删除
func (h *PermissionHandler) Delete(ctx context.Context, req *v1.PermissionDeleteRequest) (*v1.PermissionDeleteResponse, error) {
	if err := h.permissionController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call PermissionController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.PermissionDeleteResponse{}, nil
}
