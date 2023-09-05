package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type ProductHandler struct {
	controller *controller.ProductController
}

func NewProductHandler(controller *controller.ProductController) *ProductHandler {
	return &ProductHandler{controller}
}

type ProductInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

type ProductListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

type ProductListResponse []*ProductInfo

// List 产品列表
//
//	@Router			/v1/products [get]
//	@Summary		产品列表
//	@Description	产品列表
//	@Tags			产品
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string										false	"查询字符串"	format(string)
//	@Success		200		{object}	example.Success{data=ProductListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError							"服务器出错"
//	@Failure		400		{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized						"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied					"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *ProductHandler) List(ctx echo.Context) error {
	req := new(ProductListRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.ProductListRequest{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make(ProductListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &ProductInfo{
			ID:    item.ID,
			Name:  item.Name,
			Desc:  item.Desc,
			Price: item.Price,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

type ProductCreateRequest struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

// Create 产品创建
//
//	@Router			/v1/product [post]
//	@Summary		产品创建
//	@Description	产品创建
//	@Tags			产品
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ProductCreateRequest		true	"产品信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *ProductHandler) Create(ctx echo.Context) error {
	req := new(ProductCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.ProductCreateRequest{
		ProductAttr: controller.ProductAttr{
			Name:  req.Name,
			Desc:  req.Desc,
			Price: req.Price,
		},
	}
	if err := h.controller.Create(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type ProductUpdateRequest struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

// Update 产品更新
//
//	@Router			/v1/product [put]
//	@Summary		产品更新
//	@Description	产品更新
//	@Tags			产品
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ProductUpdateRequest		true	"产品信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *ProductHandler) Update(ctx echo.Context) error {
	req := new(ProductUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	p := controller.ProductUpdateRequest{
		ID: req.ID,
		ProductAttr: controller.ProductAttr{
			Name:  req.Name,
			Desc:  req.Desc,
			Price: req.Price,
		},
	}
	if err := h.controller.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type ProductDetailRequest struct {
	ID int64 `query:"id"`
}

type ProductDetailResponse = ProductInfo

// Detail 产品详情
//
//	@Router			/v1/product [get]
//	@Summary		产品详情
//	@Description	产品详情
//	@Tags			产品
//	@Accept			plain
//	@Produce		json
//	@Param			id	query		integer										true	"产品 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=ProductDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError							"服务器出错"
//	@Failure		400	{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized						"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied					"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *ProductHandler) Detail(ctx echo.Context) error {
	req := new(ProductDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &ProductDetailResponse{
		ID:    ret.ID,
		Name:  ret.Name,
		Desc:  ret.Desc,
		Price: ret.Price,
	}

	return ctx.JSON(http.StatusOK, data)
}

type ProductDeleteRequest struct {
	ID int64 `query:"id"`
}

// Delete 产品删除
//
//	@Router			/v1/product [delete]
//	@Summary		产品删除
//	@Description	产品删除
//	@Tags			产品
//	@Accept			plain
//	@Produce		json
//	@Param			id	query		integer						true	"产品 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *ProductHandler) Delete(ctx echo.Context) error {
	req := new(ProductDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
