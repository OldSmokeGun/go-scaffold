package v1

import (
	"net/http"

	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/app/domain"
	berr "go-scaffold/internal/app/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// UserHandler 用户处理器
type UserHandler struct {
	controller *controller.UserController
}

// NewUserHandler 构造用户处理器
func NewUserHandler(controller *controller.UserController) *UserHandler {
	return &UserHandler{controller}
}

// UserListRequest 用户列表请求参数
type UserListRequest struct {
	Keyword string `json:"keyword" query:"keyword"`
}

// UserListItem 用户列表项
type UserListItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// UserListResponse 用户列表响应数据
type UserListResponse []*UserListItem

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
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	p := domain.UserListParam{
		Keyword: req.Keyword,
	}
	ret, err := h.controller.List(ctx.Request().Context(), p)
	if err != nil {
		return err
	}

	data := make(UserListResponse, 0, len(ret))
	for _, item := range ret {
		data = append(data, &UserListItem{
			ID:    item.ID.Int64(),
			Name:  item.Name,
			Age:   item.Age,
			Phone: item.Phone,
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

// UserCreateRequest 创建用户请求参数
type UserCreateRequest struct {
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// Create 创建用户
//
//	@Router			/v1/user [post]
//	@Summary		创建用户
//	@Description	创建用户
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
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	p := domain.User{
		Name:  req.Name,
		Age:   req.Age,
		Phone: req.Phone,
	}
	if err := h.controller.Create(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

// UserUpdateRequest 更新用户请求参数
type UserUpdateRequest struct {
	ID    int64  `json:"id" param:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

// Update 更新用户
//
//	@Router			/v1/user/{id} [put]
//	@Summary		更新用户
//	@Description	更新用户
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			id		path		integer						true	"用户 id"	format(uint)	minimum(1)
//	@Param			data	body		UserCreateRequest			true	"用户信息"	format(string)
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
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	p := domain.User{
		ID:    domain.ID(req.ID),
		Name:  req.Name,
		Age:   req.Age,
		Phone: req.Phone,
	}
	if err := h.controller.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

// UserDetailRequest 用户详情请求参数
type UserDetailRequest struct {
	ID int64 `json:"id" param:"id"`
}

// UserDetailResponse 用户详情响应数据
type UserDetailResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
	Phone string `json:"phone"`
}

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
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	ret, err := h.controller.Detail(ctx.Request().Context(), domain.ID(req.ID))
	if err != nil {
		return err
	}

	data := &UserDetailResponse{
		ID:    ret.ID.Int64(),
		Name:  ret.Name,
		Age:   ret.Age,
		Phone: ret.Phone,
	}

	return ctx.JSON(http.StatusOK, data)
}

// UserDeleteRequest 删除用户请求参数
type UserDeleteRequest struct {
	ID int64 `json:"id" param:"id"`
}

// Delete 删除用户
//
//	@Router			/v1/user/{id} [delete]
//	@Summary		删除用户
//	@Description	删除用户
//	@Tags			用户
//	@Accept			plain
//	@Produce		json
//	@Param			id	path		integer						true	"用户 id"	format(uint)	minimum(1)
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *UserHandler) Delete(ctx echo.Context) error {
	req := new(UserDeleteRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.Wrap(err, berr.ErrBadRequest.Error())
	}

	if err := h.controller.Delete(ctx.Request().Context(), domain.ID(req.ID)); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
