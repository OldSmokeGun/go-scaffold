package v1

import (
	"context"
	"log/slog"

	v1 "go-scaffold/internal/app/adapter/server/grpc/api/v1"
	"go-scaffold/internal/app/adapter/server/grpc/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type ProductHandler struct {
	v1.UnimplementedProductServer
	logger            *slog.Logger
	productController *controller.ProductController
}

func NewProductHandler(
	logger *slog.Logger,
	productController *controller.ProductController,
) *ProductHandler {
	return &ProductHandler{
		logger:            logger,
		productController: productController,
	}
}

// List 产品列表
func (h *ProductHandler) List(ctx context.Context, req *v1.ProductListRequest) (*v1.ProductListResponse, error) {
	r := controller.ProductListRequest{
		Keyword: req.Keyword,
	}

	list, err := h.productController.List(ctx, r)
	if err != nil {
		h.logger.Error("call ProductController.List method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	items := make([]*v1.ProductInfo, 0, len(list))

	for _, item := range list {
		items = append(items, &v1.ProductInfo{
			Id:    item.ID,
			Name:  item.Name,
			Desc:  item.Desc,
			Price: int64(item.Price),
		})
	}

	return &v1.ProductListResponse{Items: items}, nil
}

// Create 产品创建
func (h *ProductHandler) Create(ctx context.Context, req *v1.ProductCreateRequest) (*v1.ProductCreateResponse, error) {
	r := controller.ProductCreateRequest{
		ProductAttr: controller.ProductAttr{
			Name:  req.Name,
			Desc:  req.Desc,
			Price: int(req.Price),
		},
	}

	if err := h.productController.Create(ctx, r); err != nil {
		h.logger.Error("call ProductController.Create method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.ProductCreateResponse{}, nil
}

// Update 产品更新
func (h *ProductHandler) Update(ctx context.Context, req *v1.ProductUpdateRequest) (*v1.ProductUpdateResponse, error) {
	r := controller.ProductUpdateRequest{
		ID: req.Id,
		ProductAttr: controller.ProductAttr{
			Name:  req.Name,
			Desc:  req.Desc,
			Price: int(req.Price),
		},
	}

	if err := h.productController.Update(ctx, r); err != nil {
		h.logger.Error("call ProductController.Update method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.ProductUpdateResponse{}, nil
}

// Detail 产品详情
func (h *ProductHandler) Detail(ctx context.Context, req *v1.ProductDetailRequest) (*v1.ProductInfo, error) {
	ret, err := h.productController.Detail(ctx, req.Id)
	if err != nil {
		h.logger.Error("call ProductController.Detail method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.ProductInfo{
		Id:    ret.ID,
		Name:  ret.Name,
		Desc:  ret.Desc,
		Price: int64(ret.Price),
	}, nil
}

// Delete 产品删除
func (h *ProductHandler) Delete(ctx context.Context, req *v1.ProductDeleteRequest) (*v1.ProductDeleteResponse, error) {
	if err := h.productController.Delete(ctx, req.Id); err != nil {
		h.logger.Error("call ProductController.Delete method error", slog.Any("error", err))
		return nil, errors.Wrap(err)
	}

	return &v1.ProductDeleteResponse{}, nil
}
