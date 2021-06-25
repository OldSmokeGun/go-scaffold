package response

const (
	SuccessCode          = "OK"                // 成功
	FailedCode           = "ERROR"             // 服务器出错
	LoginInvalidCode     = "LOGIN_VALID"       // 登陆失效
	IllegalRequestCode   = "ILLEGAL_REQUEST"   // 非法请求
	ArgumentsInvalidCode = "ARGUMENTS_INVALID" // 参数无效

	SuccessCodeMessage          = ""
	FailedCodeMessage           = "服务器出错"
	LoginInvalidCodeMessage     = "登陆失效"
	IllegalRequestCodeMessage   = "非法请求"
	ArgumentsInvalidCodeMessage = "参数格式无效"
)
