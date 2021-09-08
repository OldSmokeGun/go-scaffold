package response

const (
	SuccessCode        = 10000 // 成功
	ServerErrorCode    = 10001 // 服务器错误
	ArgumentErrorCode  = 10002 // 参数错误
	IllegalRequestCode = 10003 // 非法请求
	LoginInvalidCode   = 10004 // 登陆失效
	NoPermissionCode   = 10005 // 没有权限

	SuccessCodeMessage        = "OK"
	ServerErrorCodeMessage    = "服务器出错"
	ArgumentErrorCodeMessage  = "参数无效"
	IllegalRequestCodeMessage = "非法请求"
	LoginInvalidCodeMessage   = "登陆失效"
	NoPermissionCodeMessage   = "暂无权限"
)

var (
	Success        = NewResponse(SuccessCode, SuccessCodeMessage, nil)
	ServerError    = NewResponse(ServerErrorCode, ServerErrorCodeMessage, nil)
	ArgumentError  = NewResponse(ArgumentErrorCode, ArgumentErrorCodeMessage, nil)
	IllegalRequest = NewResponse(IllegalRequestCode, IllegalRequestCodeMessage, nil)
	LoginInvalid   = NewResponse(LoginInvalidCode, LoginInvalidCodeMessage, nil)
	NoPermission   = NewResponse(NoPermissionCode, NoPermissionCodeMessage, nil)
)
