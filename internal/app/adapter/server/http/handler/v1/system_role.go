package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	httperr "go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type SystemRoleHandler struct {
	controller *controller.SystemRoleController
}

func NewSystemRoleHandler(controller *controller.SystemRoleController) *SystemRoleHandler {
	return &SystemRoleHandler{controller}
}

type SystemRoleInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SystemRoleListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

type SystemRoleListResponse []*SystemRoleInfo

// List 角色列表
//
//	@Router			/v1/system-roles [get]
//	@Summary		列表
//	@Description	角色列表
//	@Tags			角色
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string									false	"查询字符串"	format(string)
//	@Success		200		{object}	example.Success{data=SystemRoleListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError						"服务器出错"
//	@Failure		400		{object}	example.ClientError						"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized					"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied				"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound				"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest					"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) List(ctx echo.Context) error {
	req := new(SystemRoleListRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemRoleListRequest{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make(SystemRoleListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &SystemRoleInfo{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemRoleCreateRequest struct {
	Name string `json:"name"`
}

// Create 角色创建
//
//	@Router			/v1/system-role [post]
//	@Summary		角色创建
//	@Description	角色创建
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemRoleCreateRequest			true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) Create(ctx echo.Context) error {
	req := new(SystemRoleCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemRoleCreateRequest{
		SystemRoleAttr: controller.SystemRoleAttr{
			Name: req.Name,
		},
	}
	if err := h.controller.Create(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemRoleUpdateRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Update 角色更新
//
//	@Router			/v1/system-role [put]
//	@Summary		角色更新
//	@Description	角色更新
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemRoleUpdateRequest			true	"权限信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) Update(ctx echo.Context) error {
	req := new(SystemRoleUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	p := controller.SystemRoleUpdateRequest{
		ID: req.ID,
		SystemRoleAttr: controller.SystemRoleAttr{
			Name: req.Name,
		},
	}
	if err := h.controller.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemRoleDetailRequest struct {
	ID int64 `param:"id"`
}

type SystemRoleDetailResponse = SystemRoleInfo

// Detail 角色详情
//
//	@Router			/v1/system-role/{id} [get]
//	@Summary		角色详情
//	@Description	角色详情
//	@Tags			角色
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer										true	"角色 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=SystemRoleDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError							"服务器出错"
//	@Failure		400	{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized						"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied					"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) Detail(ctx echo.Context) error {
	req := new(SystemRoleDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &SystemRoleDetailResponse{
		ID:   ret.ID,
		Name: ret.Name,
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemRoleDeleteRequest struct {
	ID int64 `param:"id"`
}

// Delete 角色删除
//
//	@Router			/v1/system-role/{id} [delete]
//	@Summary		角色删除
//	@Description	角色删除
//	@Tags			角色
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer						true	"角色 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) Delete(ctx echo.Context) error {
	req := new(SystemRoleDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemRoleGrantPermissionsRequest struct {
	Role        int64   `json:"role"`
	Permissions []int64 `json:"permissions"`
}

// GrantPermissions 授权角色权限
//
//	@Router			/v1/system-role/permissions [post]
//	@Summary		授权角色权限
//	@Description	授权角色权限
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemRoleGrantPermissionsRequest	true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) GrantPermissions(ctx echo.Context) error {
	req := new(SystemRoleGrantPermissionsRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemRoleGrantPermissionsRequest{
		Role:        req.Role,
		Permissions: req.Permissions,
	}
	if err := h.controller.GrantPermissions(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemRoleGetPermissionsRequest struct {
	ID int64 `query:"id"`
}

type SystemRoleGetPermissionsResponse []*SystemPermissionInfo

// GetPermissions 获取角色权限
//
//	@Router			/v1/system-role/permissions [get]
//	@Summary		获取角色权限
//	@Description	获取角色权限
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemRoleGetPermissionsRequest	true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemRoleHandler) GetPermissions(ctx echo.Context) error {
	req := new(SystemRoleGetPermissionsRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.GetPermissions(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := make(SystemRoleGetPermissionsResponse, 0, len(ret))
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
