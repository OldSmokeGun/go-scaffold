package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-scaffold/internal/app/adapter/server/http/middleware"
	httperr "go-scaffold/internal/app/adapter/server/http/pkg/errors"
	"go-scaffold/internal/app/controller"
	"go-scaffold/internal/app/domain"
)

type SystemSessionHandler struct {
	controller *controller.SystemSessionController
}

func NewSystemSessionHandler(controller *controller.SystemSessionController) *SystemSessionHandler {
	return &SystemSessionHandler{controller}
}

type SystemSessionRegisterRequest SystemUserCreateRequest

type SystemSessionRegisterResponse struct {
	User  *SystemUserInfo `json:"user"`
	Token string          `json:"token"`
}

// Register 注册
//
//	@Router			/v1/register [post]
//	@Summary		注册
//	@Description	注册
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemSessionRegisterRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=SystemSessionRegisterResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError								"服务器出错"
//	@Failure		400		{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized							"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied						"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) Register(ctx echo.Context) error {
	req := new(SystemSessionRegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemSessionRegisterRequest{
		SystemUserAttr: controller.SystemUserAttr{
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

	data := SystemSessionRegisterResponse{
		User: &SystemUserInfo{
			ID:       ret.User.ID,
			Username: ret.User.Username,
			Nickname: ret.User.Nickname,
			Phone:    ret.User.Phone,
		},
		Token: ret.Token,
	}

	return ctx.JSON(http.StatusOK, data)
}

type SystemSessionLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SystemSessionLoginResponse struct {
	User  *SystemSessionProfileResponse `json:"user"`
	Token string                        `json:"token"`
}

// Login 登录
//
//	@Router			/v1/login [post]
//	@Summary		登录
//	@Description	登录
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemSessionLoginRequest							true	"请求体"	format(string)
//	@Success		200		{object}	example.Success{data=SystemSessionLoginResponse}	"成功响应"
//	@Failure		500		{object}	example.ServerError							"服务器出错"
//	@Failure		400		{object}	example.ClientError							"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized						"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied					"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound					"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest						"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) Login(ctx echo.Context) error {
	req := new(SystemSessionLoginRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	r := controller.SystemSessionLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	ret, err := h.controller.Login(ctx.Request().Context(), r)
	if err != nil {
		return err
	}

	data := SystemSessionLoginResponse{
		User:  newSystemSessionProfile(ret.User),
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
func (h *SystemSessionHandler) Logout(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	if err := h.controller.Logout(ctx.Request().Context(), user.ID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemSessionUpdateProfileRequest struct {
	Nickname string `json:"nickname"`
}

// UpdateProfile 更新账号信息
//
//	@Router			/v1/profile [put]
//	@Summary		更新账号信息
//	@Description	更新账号信息
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemSessionUpdateProfileRequest	true	"请求体"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) UpdateProfile(ctx echo.Context) error {
	req := new(SystemSessionUpdateProfileRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	user := ctx.(*middleware.Context).GetUser()

	r := controller.SystemSessionUpdateProfileRequest{
		ID:       user.ID,
		Nickname: req.Nickname,
	}
	if err := h.controller.UpdateProfile(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemSessionUpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

// UpdatePassword 修改当前登录用户密码
//
//	@Router			/v1/profile/password [put]
//	@Summary		修改当前登录用户密码
//	@Description	修改当前登录用户密码。新密码 8～18 位，且同时包含数字、大写字母、小写字母和符号
//	@Tags			账号
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SystemSessionUpdatePasswordRequest	true	"原密码与新密码"	format(string)
//	@Success		200		{object}	example.Success				"成功响应"
//	@Failure		500		{object}	example.ServerError			"服务器出错"
//	@Failure		400		{object}	example.ClientError			"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401		{object}	example.Unauthorized		"登陆失效"
//	@Failure		403		{object}	example.PermissionDenied	"没有权限"
//	@Failure		404		{object}	example.ResourceNotFound	"资源不存在"
//	@Failure		429		{object}	example.TooManyRequest		"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) UpdatePassword(ctx echo.Context) error {
	req := new(SystemSessionUpdatePasswordRequest)
	if err := ctx.Bind(req); err != nil {
		return httperr.WrapHTTTPError(err).SetMessage("request parameter parsing error")
	}

	user := ctx.(*middleware.Context).GetUser()
	r := controller.SystemSessionUpdatePasswordRequest{
		ID:          user.ID,
		OldPassword: req.OldPassword,
		Password:    req.Password,
	}
	if err := h.controller.UpdatePassword(ctx.Request().Context(), r); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type SystemSessionProfileResponse struct {
	ID       int64             `json:"id"`
	Username string            `json:"username"`
	Nickname string            `json:"nickname"`
	Phone    string            `json:"phone"`
	Roles    []*SystemRoleInfo `json:"roles"`
}

func newSystemSessionProfile(user *domain.SystemUserProfile) *SystemSessionProfileResponse {
	if user == nil {
		return nil
	}
	roles := make([]*SystemRoleInfo, 0, len(user.Roles))
	for _, role := range user.Roles {
		if role == nil {
			continue
		}
		roles = append(roles, &SystemRoleInfo{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return &SystemSessionProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Phone:    user.Phone,
		Roles:    roles,
	}
}

// GetProfile 获取账号信息
//
//	@Router			/v1/profile [get]
//	@Summary		获取账号信息
//	@Description	获取账号信息
//	@Tags			账号
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	example.Success{data=SystemSessionProfileResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError								"服务器出错"
//	@Failure		400	{object}	example.ClientError								"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized							"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied						"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound						"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest							"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) GetProfile(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	ret, err := h.controller.GetProfile(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, newSystemSessionProfile(ret))
}

type SystemSessionGetPermissionsResponse []*SystemPermissionInfo

// GetPermissions 获取账号权限
//
//	@Router			/v1/permissions [get]
//	@Summary		获取账号权限
//	@Description	获取账号权限
//	@Tags			账号
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	example.Success{data=SystemSessionGetPermissionsResponse}	"成功响应"
//	@Failure		500	{object}	example.ServerError									"服务器出错"
//	@Failure		400	{object}	example.ClientError									"客户端请求错误（code 类型应为 int，string 仅为了表达多个错误码）"
//	@Failure		401	{object}	example.Unauthorized								"登陆失效"
//	@Failure		403	{object}	example.PermissionDenied							"没有权限"
//	@Failure		404	{object}	example.ResourceNotFound							"资源不存在"
//	@Failure		429	{object}	example.TooManyRequest								"请求过于频繁"
//	@Security		Authorization
func (h *SystemSessionHandler) GetPermissions(ctx echo.Context) error {
	user := ctx.(*middleware.Context).GetUser()

	ret, err := h.controller.GetPermissions(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	data := make(SystemSessionGetPermissionsResponse, 0, len(ret))
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
