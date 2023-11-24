package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type UserHandler struct {
	controller *controller.UserController
}

func NewUserHandler(controller *controller.UserController) *UserHandler {
	return &UserHandler{controller}
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type UserListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

type UserListResponse []*UserInfo

// List 用户列表
//
//	@Router			/v1/users [get]
//	@Summary		用户列表
//	@Description	用户列表
//	@Tags			用户
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			keyword	query		string									false	"查询字符串"	format(string)
//	@Success		200		{object}	example.Success{data=UserListResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError						"服务器出错"
//	@Failure		400		{object}	example.ClientError						"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized					"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied				"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound				"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest					"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) List(ctx echo.Context) error {
	req := new(UserListRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.UserListRequest{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := make(UserListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &UserInfo{
			ID:       item.ID,
			Username: item.Username,
			Nickname: item.Nickname,
			Phone:    item.Phone,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// Create 用户创建
//
//	@Router			/v1/user [post]
//	@Summary		用户创建
//	@Description	用户创建
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UserCreateRequest			true	"用户信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) Create(ctx echo.Context) error {
	req := new(UserCreateRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.UserCreateRequest{
		UserAttr: controller.UserAttr{
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

type UserUpdateRequest struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

// Update 用户更新
//
//	@Router			/v1/user [put]
//	@Summary		用户更新
//	@Description	用户更新
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UserUpdateRequest			true	"用户信息"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) Update(ctx echo.Context) error {
	req := new(UserUpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	p := controller.UserUpdateRequest{
		ID: req.ID,
		UserAttr: controller.UserAttr{
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

type UserDetailRequest struct {
	ID int64 `param:"id"`
}

type UserDetailResponse = UserInfo

// Detail 用户详情
//
//	@Router			/v1/user/{id} [get]
//	@Summary		用户详情
//	@Description	用户详情
//	@Tags			用户
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer										true	"用户 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success{data=UserDetailResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError							"服务器出错"
//	@Failure		400	{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized						"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied					"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) Detail(ctx echo.Context) error {
	req := new(UserDetailRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := &UserDetailResponse{
		ID:       ret.ID,
		Username: ret.Username,
		Nickname: ret.Nickname,
		Phone:    ret.Phone,
	}

	return ctx.JSON(http.StatusOK, data)
}

type UserDeleteRequest struct {
	ID int64 `param:"id"`
}

// Delete 用户删除
//
//	@Router			/v1/user/{id} [delete]
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
func (h *UserHandler) Delete(ctx echo.Context) error {
	req := new(PermissionDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	if err := h.controller.Delete(ctx.Request().Context(), req.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type UserAssignRoleRequest struct {
	User  int64   `json:"user"`
	Roles []int64 `json:"roles"`
}

// AssignRoles 分配用户角色
//
//	@Router			/v1/user/roles [post]
//	@Summary		分配用户角色
//	@Description	分配用户角色
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UserAssignRoleRequest		true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) AssignRoles(ctx echo.Context) error {
	req := new(UserAssignRoleRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.UserAssignRoleRequest{
		User:  req.User,
		Roles: req.Roles,
	}
	if err := h.controller.AssignRoles(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type UserGetRoleRequest struct {
	ID int64 `query:"id"`
}

type UserGetRoleResponse []*RoleInfo

// GetRoles 获取用户角色
//
//	@Router			/v1/user/roles [get]
//	@Summary		获取用户角色
//	@Description	获取用户角色
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UserGetRoleRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=UserGetRoleResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError							"服务器出错"
//	@Failure		400		{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized						"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied					"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) GetRoles(ctx echo.Context) error {
	req := new(UserGetRoleRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	ret, err := h.controller.GetRoles(ctx.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	data := make(UserGetRoleResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &RoleInfo{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}
