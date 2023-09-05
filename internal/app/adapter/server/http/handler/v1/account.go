package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/adapter/server/http/middleware"
	"go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
)

type AccountHandler struct {
	controller *controller.AccountController
}

func NewAccountHandler(controller *controller.AccountController) *AccountHandler {
	return &AccountHandler{controller}
}

type AccountRegisterRequest UserCreateRequest

type AccountRegisterResponse struct {
	User  *UserInfo `json:"user"`
	Token string    `json:"token"`
}

// Register 注册
//
//	@Router			/v1/register [post]
//	@Summary		注册
//	@Description	注册
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		AccountRegisterRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=AccountRegisterResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError								"服务器出错"
//	@Failure		400		{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized							"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied						"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) Register(ctx echo.Context) error {
	req := new(AccountRegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.AccountRegisterRequest{
		UserAttr: controller.UserAttr{
			Username: req.Username,
			Password: req.Password,
			Nickname: req.Nickname,
			Phone:    req.Phone,
		},
	}
	ret, err := h.controller.Register(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := AccountRegisterResponse{
		User: &UserInfo{
			ID:       ret.User.ID,
			Username: ret.User.Username,
			Nickname: ret.User.Nickname,
			Phone:    ret.User.Phone,
		},
		Token: ret.Token,
	}

	return ctx.JSON(http.StatusOK, data)
}

type AccountLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountLoginResponse struct {
	User  *UserInfo `json:"user"`
	Token string    `json:"token"`
}

// Login 登录
//
//	@Router			/v1/login [post]
//	@Summary		登录
//	@Description	登录
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		AccountLoginRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=AccountLoginResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError							"服务器出错"
//	@Failure		400		{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized						"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied					"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) Login(ctx echo.Context) error {
	req := new(AccountLoginRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	r := controller.AccountLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	ret, err := h.controller.Login(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := AccountLoginResponse{
		User: &UserInfo{
			ID:       ret.User.ID,
			Username: ret.User.Username,
			Nickname: ret.User.Nickname,
			Phone:    ret.User.Phone,
		},
		Token: ret.Token,
	}

	return ctx.JSON(http.StatusOK, data)
}

// Logout 登出
//
//	@Router			/v1/logout [delete]
//	@Summary		登出
//	@Description	登出
//	@Tags			账号
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	example.Success				"成功响应"
//	@Failure		500	{object}	example.ServerError			"服务器出错"
//	@Failure		400	{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized		"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied	"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) Logout(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	if err := h.controller.Logout(ctx.Request().Context(), user.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type AccountUpdateProfileRequest struct {
	Nickname string `json:"nickname"`
}

// UpdateProfile 更新账号信息
//
//	@Router			/v1/account/profile [put]
//	@Summary		更新账号信息
//	@Description	更新账号信息
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		AccountUpdateProfileRequest	true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) UpdateProfile(ctx echo.Context) error {
	req := new(AccountUpdateProfileRequest)
	if err := ctx.Bind(req); err != nil {
		return errors.WrapHTTTPError(err.(*echo.HTTPError)).SetMessage("request parameter parsing error").Unwrap()
	}

	user := ctx.(*middleware.Context).GetUser()

	r := controller.AccountUpdateProfileRequest{
		ID:       user.ID,
		Nickname: req.Nickname,
	}
	if err := h.controller.UpdateProfile(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type AccountProfileResponse UserInfo

// GetProfile 获取账号信息
//
//	@Router			/v1/account/profile [get]
//	@Summary		获取账号信息
//	@Description	获取账号信息
//	@Tags			账号
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	example.Success{data=AccountProfileResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError								"服务器出错"
//	@Failure		400	{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized							"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied						"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) GetProfile(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	ret, err := h.controller.GetProfile(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	data := &AccountProfileResponse{
		ID:       ret.ID,
		Username: ret.Username,
		Nickname: ret.Nickname,
		Phone:    ret.Phone,
	}

	return ctx.JSON(http.StatusOK, data)
}

type AccountGetPermissionsResponse []*PermissionInfo

// GetPermissions 获取账号权限
//
//	@Router			/v1/account/permissions [get]
//	@Summary		获取账号权限
//	@Description	获取账号权限
//	@Tags			账号
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	example.Success{data=AccountGetPermissionsResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError									"服务器出错"
//	@Failure		400	{object}	example.ClientError									"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized								"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied							"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound							"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest								"请求过于频繁"
//	@Security		Authorization
func (h *AccountHandler) GetPermissions(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	ret, err := h.controller.GetPermissions(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	data := make(AccountGetPermissionsResponse, 0, len(ret))
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
