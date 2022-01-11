package responsex

// StatusCode 响应状态码
type StatusCode int

const (
	SuccessCode          StatusCode = 10000 // 成功
	ServerErrorCode      StatusCode = 10001 // 服务器出错
	ClientErrorCode      StatusCode = 10002 // 客户端请求错误
	ValidateErrorCode    StatusCode = 10003 // 参数校验错误
	UnauthorizedCode     StatusCode = 10004 // 登陆失效
	PermissionDeniedCode StatusCode = 10005 // 没有权限
	ResourceNotFoundCode StatusCode = 10006 // 资源不存在
	TooManyRequestCode   StatusCode = 10007 // 请求过于频繁
)

func (r StatusCode) String() string {
	switch r {
	case SuccessCode:
		return "OK"
	case ServerErrorCode:
		return "服务器出错"
	case ClientErrorCode:
		return "客户端请求错误"
	case ValidateErrorCode:
		return "参数校验错误"
	case UnauthorizedCode:
		return "登陆失效"
	case PermissionDeniedCode:
		return "暂无权限"
	case ResourceNotFoundCode:
		return "资源不存在"
	case TooManyRequestCode:
		return "请求过于频繁"
	}

	return ""
}

var (
	Success          = NewResponse(SuccessCode, SuccessCode.String(), nil)
	ServerError      = NewResponse(ServerErrorCode, ServerErrorCode.String(), nil)
	ClientError      = NewResponse(ClientErrorCode, ClientErrorCode.String(), nil)
	ValidateError    = NewResponse(ValidateErrorCode, ValidateErrorCode.String(), nil)
	Unauthorized     = NewResponse(UnauthorizedCode, UnauthorizedCode.String(), nil)
	PermissionDenied = NewResponse(PermissionDeniedCode, PermissionDeniedCode.String(), nil)
)
