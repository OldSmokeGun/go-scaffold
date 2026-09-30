package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type SystemPermissionHandler struct {
	v1.UnimplementedSystemPermissionServer
	logger               *slog.Logger
	permissionController *controller.SystemPermissionController
}

func NewSystemPermissionHandler(
	logger *slog.Logger,
	permissionController *controller.SystemPermissionController,
) *SystemPermissionHandler {
	return &SystemPermissionHandler{
		logger:               logger,
		permissionController: permissionController,
	}
}

// List 权限列表
func (h *SystemPermissionHandler) List(ctx context.Context, req *v1.SystemPermissionListRequest) (*v1.SystemPermissionListResponse, error) {
	r := controller.SystemPermissionListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.permissionController.List(ctx, r)
	if err != nil {
		h.logger.Error("call SystemPermissionController.List method error", slog.Any("error", err))
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

	return &v1.SystemPermissionListResponse{Items: items}, nil
}

// Create 权限创建
func (h *SystemPermissionHandler) Create(ctx context.Context, req *v1.SystemPermissionCreateRequest) (*v1.SystemPermissionCreateResponse, error) {
	r := controller.SystemPermissionCreateRequest{
		SystemPermissionAttr: controller.SystemPermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}

	if err := h.permissionController.Create(ctx, r); err != nil {
		h.logger.Error("call SystemPermissionController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemPermissionCreateResponse{}, nil
}

// Update 权限更新
func (h *SystemPermissionHandler) Update(ctx context.Context, req *v1.SystemPermissionUpdateRequest) (*v1.SystemPermissionUpdateResponse, error) {
	r := controller.SystemPermissionUpdateRequest{
		ID: req.Id,
		SystemPermissionAttr: controller.SystemPermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}

	if err := h.permissionController.Update(ctx, r); err != nil {
		h.logger.Error("call SystemPermissionController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemPermissionUpdateResponse{}, nil
}

// Detail 权限详情
func (h *SystemPermissionHandler) Detail(ctx context.Context, req *v1.SystemPermissionDetailRequest) (*v1.SystemPermissionInfo, error) {
	ret, err := h.permissionController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call SystemPermissionController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemPermissionInfo{
		Id:       ret.ID,
		Key:      ret.Key,
		Name:     ret.Name,
		Desc:     ret.Desc,
		ParentID: ret.ParentID,
	}, nil
}

// Delete 权限删除
func (h *SystemPermissionHandler) Delete(ctx context.Context, req *v1.SystemPermissionDeleteRequest) (*v1.SystemPermissionDeleteResponse, error) {
	if err := h.permissionController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call SystemPermissionController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.SystemPermissionDeleteResponse{}, nil
}
