package example

// Success 成功
type Success struct {
	Data any `json:"data,omitempty"`
}

// ServerError 服务器出错
type ServerError struct {
	ErrNo  int    `json:"errNo" example:"10001"`
	ErrMsg string `json:"errMsg" example:"服务器出错"`
}

// ClientError 客户端请求错误
type ClientError struct {
	ErrNo  string `json:"errNo" example:"10002|10003"` // errNo 类型应为 int，string 仅为了表达多个错误码
	ErrMsg string `json:"errMsg" example:"客户端请求错误|参数校验错误"`
}

// Unauthorized 未经授权
type Unauthorized struct {
	ErrNo  int    `json:"errNo" example:"10004"`
	ErrMsg string `json:"errMsg" example:"未经授权"`
}

// PermissionDenied 没有权限
type PermissionDenied struct {
	ErrNo  int    `json:"errNo" example:"10005"`
	ErrMsg string `json:"errMsg" example:"暂无权限"`
}

// ResourceNotFound 资源不存在
type ResourceNotFound struct {
	ErrNo  int    `json:"errNo" example:"10006"`
	ErrMsg string `json:"errMsg" example:"资源不存在"`
}

// TooManyRequest 请求过于频繁
type TooManyRequest struct {
	ErrNo  int    `json:"errNo" example:"10007"`
	ErrMsg string `json:"errMsg" example:"请求过于频繁"`
}
