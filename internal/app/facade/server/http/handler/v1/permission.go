package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/controller"
	httperr "go-scaffold/internal/app/facade/server/http/pkg/errors"
)

type PermissionHandler struct {
	controller *controller.PermissionController
}

func NewPermissionHandler(controller *controller.PermissionController) *PermissionHandler {
	return &PermissionHandler{controller}
}

type PermissionInfo struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

type PermissionListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

type PermissionListResponse []*PermissionInfo

// List 权限列表
//
//	@Router			/v1/permissions [get]
//	@Summary		权限列表
//	@Description	权限列表
//	@Tags			权限
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string											false	"查询字符串"	format(string)
//	@Success		200		{object}	example.Success{data=PermissionListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError								"服务器出错"
//	@Failure		400		{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized							"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied						"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *PermissionHandler) List(ctx echo.Context) error {
	req := new(PermissionListRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	r := controller.PermissionListRequest{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make(PermissionListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &PermissionInfo{
			ID:       item.ID,
			Key:      item.Key,
			Name:     item.Name,
			Desc:     item.Desc,
			ParentID: item.ParentID,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

type PermissionCreateRequest struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

// Create 权限创建
//
//	@Router			/v1/permission [post]
//	@Summary		权限创建
//	@Description	权限创建
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Param			data	body		PermissionCreateRequest		true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *PermissionHandler) Create(ctx echo.Context) error {
	req := new(PermissionCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	r := controller.PermissionCreateRequest{
		PermissionAttr: controller.PermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}
	if err := h.controller.Create(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type PermissionUpdateRequest struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

// Update 权限更新
//
//	@Router			/v1/permission [put]
//	@Summary		权限更新
//	@Description	权限更新
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Param			data	body		PermissionUpdateRequest		true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *PermissionHandler) Update(ctx echo.Context) error {
	req := new(PermissionUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	p := controller.PermissionUpdateRequest{
		ID: req.ID,
		PermissionAttr: controller.PermissionAttr{
			Key:      req.Key,
			Name:     req.Name,
			Desc:     req.Desc,
			ParentID: req.ParentID,
		},
	}
	if err := h.controller.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type PermissionDetailRequest struct {
	ID int64 `param:"id"`
}

type PermissionDetailResponse = PermissionInfo

// Detail 权限详情
//
//	@Router			/v1/permission/{id} [get]
//	@Summary		权限详情
//	@Description	权限详情
//	@Tags			权限
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer											true	"权限 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=PermissionDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError								"服务器出错"
//	@Failure		400	{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized							"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied						"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *PermissionHandler) Detail(ctx echo.Context) error {
	req := new(PermissionDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &PermissionDetailResponse{
		ID:       ret.ID,
		Key:      ret.Key,
		Name:     ret.Name,
		Desc:     ret.Desc,
		ParentID: ret.ParentID,
	}

	return ctx.JSON(http.StatusOK, data)
}

type PermissionDeleteRequest struct {
	ID int64 `param:"id"`
}

// Delete 权限删除
//
//	@Router			/v1/permission/{id} [delete]
//	@Summary		权限删除
//	@Description	权限删除
//	@Tags			权限
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer						true	"权限 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *PermissionHandler) Delete(ctx echo.Context) error {
	req := new(PermissionDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error")
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
