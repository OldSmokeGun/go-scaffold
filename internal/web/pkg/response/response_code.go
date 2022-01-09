package response

const (
	SuccessCode          = 10000 // 成功
	ServerErrorCode      = 10001 // 服务器错误
	ArgumentErrorCode    = 10002 // 参数错误
	LoginInvalidCode     = 10003 // 登陆失效
	PermissionDeniedCode = 10004 // 没有权限

	SuccessCodeMessage          = "OK"
	ServerErrorCodeMessage      = "服务器出错"
	ArgumentErrorCodeMessage    = "参数错误"
	LoginInvalidCodeMessage     = "登陆失效"
	PermissionDeniedCodeMessage = "暂无权限"
)

var (
	Success          = NewResponse(SuccessCode, SuccessCodeMessage, nil)
	ServerError      = NewResponse(ServerErrorCode, ServerErrorCodeMessage, nil)
	ArgumentError    = NewResponse(ArgumentErrorCode, ArgumentErrorCodeMessage, nil)
	LoginInvalid     = NewResponse(LoginInvalidCode, LoginInvalidCodeMessage, nil)
	PermissionDenied = NewResponse(PermissionDeniedCode, PermissionDeniedCodeMessage, nil)
)
