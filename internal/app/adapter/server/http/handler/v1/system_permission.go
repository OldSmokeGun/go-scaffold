package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	httperr "go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type SystemPermissionHandler struct {
	controller *controller.SystemPermissionController
}

func NewSystemPermissionHandler(controller *controller.SystemPermissionController) *SystemPermissionHandler {
	return &SystemPermissionHandler{controller}
}

type SystemPermissionInfo struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

type SystemPermissionListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

type SystemPermissionListResponse []*SystemPermissionInfo

// List 权限列表
//
//	@Router			/v1/system-permissions [get]
//	@Summary		权限列表
//	@Description	权限列表
//	@Tags			权限
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string											false	"查询字符串"	format(string)
//	@Success		200		{object}	example.Success{data=SystemPermissionListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError								"服务器出错"
//	@Failure		400		{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized							"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied						"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *SystemPermissionHandler) List(ctx echo.Context) error {
	req := new(SystemPermissionListRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemPermissionListRequest{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make(SystemPermissionListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &SystemPermissionInfo{
			ID:       item.ID,
			Key:      item.Key,
			Name:     item.Name,
			Desc:     item.Desc,
			ParentID: item.ParentID,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemPermissionCreateRequest struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

// Create 权限创建
//
//	@Router			/v1/system-permission [post]
//	@Summary		权限创建
//	@Description	权限创建
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemPermissionCreateRequest		true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemPermissionHandler) Create(ctx echo.Context) error {
	req := new(SystemPermissionCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemPermissionCreateRequest{
		SystemPermissionAttr: controller.SystemPermissionAttr{
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

type SystemPermissionUpdateRequest struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}

// Update 权限更新
//
//	@Router			/v1/system-permission [put]
//	@Summary		权限更新
//	@Description	权限更新
//	@Tags			权限
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemPermissionUpdateRequest		true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemPermissionHandler) Update(ctx echo.Context) error {
	req := new(SystemPermissionUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	p := controller.SystemPermissionUpdateRequest{
		ID: req.ID,
		SystemPermissionAttr: controller.SystemPermissionAttr{
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

type SystemPermissionDetailRequest struct {
	ID int64 `param:"id"`
}

type SystemPermissionDetailResponse = SystemPermissionInfo

// Detail 权限详情
//
//	@Router			/v1/system-permission/{id} [get]
//	@Summary		权限详情
//	@Description	权限详情
//	@Tags			权限
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer											true	"权限 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=SystemPermissionDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError								"服务器出错"
//	@Failure		400	{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized							"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied						"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *SystemPermissionHandler) Detail(ctx echo.Context) error {
	req := new(SystemPermissionDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &SystemPermissionDetailResponse{
		ID:       ret.ID,
		Key:      ret.Key,
		Name:     ret.Name,
		Desc:     ret.Desc,
		ParentID: ret.ParentID,
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemPermissionDeleteRequest struct {
	ID int64 `param:"id"`
}

// Delete 权限删除
//
//	@Router			/v1/system-permission/{id} [delete]
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
func (h *SystemPermissionHandler) Delete(ctx echo.Context) error {
	req := new(SystemPermissionDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
