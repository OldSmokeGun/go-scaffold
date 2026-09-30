package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	httperr "go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/pkg/pagenation"
)

type SystemUserHandler struct {
	controller *controller.SystemUserController
}

func NewSystemUserHandler(controller *controller.SystemUserController) *SystemUserHandler {
	return &SystemUserHandler{controller}
}

type SystemUserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type SystemUserListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
	pagenation.Param
}

type SystemUserListItem struct {
	SystemUserInfo
	Roles []*SystemRoleInfo `json:"roles"`
}

type SystemUserListResponse struct {
	List  []*SystemUserListItem `json:"list"`
	Total int64                 `json:"total"`
}

// List 用户列表
//
//	@Router			/v1/system-users [get]
//	@Summary		用户列表
//	@Description	用户列表
//	@Tags			用户
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string									false	"查询字符串"	format(string)
//	@Param			page	query		int										false	"页码，从 1 开始"
//	@Param			limit	query		int										false	"每页条数，默认 10；-1 表示不限制"
//	@Success		200		{object}	example.Success{data=SystemUserListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError						"服务器出错"
//	@Failure		400		{object}	example.ClientError						"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized					"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied				"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound				"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest					"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) List(ctx echo.Context) error {
	req := new(SystemUserListRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemUserListRequest{
		Keyword: req.Keyword,
		Param:   req.Param,
	}
	ret, total, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make([]*SystemUserListItem, 0, len(ret))
	for _, item := range ret {
		roles := make([]*SystemRoleInfo, 0, len(item.Roles))
		for _, role := range item.Roles {
			roles = append(roles, &SystemRoleInfo{
				ID:   role.ID,
				Name: role.Name,
			})
		}
		data = append(data, &SystemUserListItem{
			SystemUserInfo: SystemUserInfo{
				ID:       item.ID,
				Username: item.Username,
				Nickname: item.Nickname,
				Phone:    item.Phone,
			},
			Roles: roles,
		})
	}

	return ctx.JSON(http.StatusOK, SystemUserListResponse{List: data, Total: total})
}

type SystemUserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// Create 用户创建
//
//	@Router			/v1/system-user [post]
//	@Summary		用户创建
//	@Description	用户创建
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemUserCreateRequest			true	"用户信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) Create(ctx echo.Context) error {
	req := new(SystemUserCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemUserCreateRequest{
		SystemUserAttr: controller.SystemUserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}
	if err := h.controller.Create(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemUserUpdateRequest struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// Update 用户更新
//
//	@Router			/v1/system-user [put]
//	@Summary		用户更新
//	@Description	用户更新
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemUserUpdateRequest			true	"用户信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) Update(ctx echo.Context) error {
	req := new(SystemUserUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	p := controller.SystemUserUpdateRequest{
		ID: req.ID,
		SystemUserAttr: controller.SystemUserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}
	if err := h.controller.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemUserDetailRequest struct {
	ID int64 `param:"id"`
}

type SystemUserDetailResponse = SystemUserInfo

// Detail 用户详情
//
//	@Router			/v1/system-user/{id} [get]
//	@Summary		用户详情
//	@Description	用户详情
//	@Tags			用户
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer										true	"用户 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=SystemUserDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError							"服务器出错"
//	@Failure		400	{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized						"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied					"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) Detail(ctx echo.Context) error {
	req := new(SystemUserDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &SystemUserDetailResponse{
		ID:       ret.ID,
		Username: ret.Username,
		Nickname: ret.Nickname,
		Phone:    ret.Phone,
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemUserDeleteRequest struct {
	ID int64 `param:"id"`
}

// Delete 用户删除
//
//	@Router			/v1/system-user/{id} [delete]
//	@Summary		用户删除
//	@Description	用户删除
//	@Tags			用户
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
func (h *SystemUserHandler) Delete(ctx echo.Context) error {
	req := new(SystemPermissionDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemUserAssignRoleRequest struct {
	User  int64   `json:"user"`
	Roles []int64 `json:"roles"`
}

// AssignRoles 分配用户角色
//
//	@Router			/v1/system-user/roles [post]
//	@Summary		分配用户角色
//	@Description	分配用户角色
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemUserAssignRoleRequest		true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) AssignRoles(ctx echo.Context) error {
	req := new(SystemUserAssignRoleRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemUserAssignRoleRequest{
		User:  req.User,
		Roles: req.Roles,
	}
	if err := h.controller.AssignRoles(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemUserGetRoleRequest struct {
	ID int64 `query:"id"`
}

type SystemUserGetRoleResponse []*SystemRoleInfo

// GetRoles 获取用户角色
//
//	@Router			/v1/system-user/roles [get]
//	@Summary		获取用户角色
//	@Description	获取用户角色
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemUserGetRoleRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=SystemUserGetRoleResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError							"服务器出错"
//	@Failure		400		{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized						"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied					"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *SystemUserHandler) GetRoles(ctx echo.Context) error {
	req := new(SystemUserGetRoleRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	ret, err := h.controller.GetRoles(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := make(SystemUserGetRoleResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &SystemRoleInfo{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}
