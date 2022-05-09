package apiexample

// Success 成功
type Success struct {
	Code int         `json:"code" example:"10000"`
	Msg  string      `json:"msg" example:"OK"`
	Data interface{} `json:"data,omitempty"`
}

// ServerError 服务器出错
type ServerError struct {
	Code int    `json:"code" example:"10001"`
	Msg  string `json:"msg" example:"服务器出错"`
}

// ClientError 客户端请求错误
type ClientError struct {
	Code string `json:"code" example:"10002|10003"` // code 类型应为 int，string 仅为了表达多个错误码
	Msg  string `json:"msg" example:"客户端请求错误|参数校验错误"`
}

// Unauthorized 未经授权
type Unauthorized struct {
	Code int    `json:"code" example:"10004"`
	Msg  string `json:"msg" example:"未经授权"`
}

// PermissionDenied 没有权限
type PermissionDenied struct {
	Code int    `json:"code" example:"10005"`
	Msg  string `json:"msg" example:"暂无权限"`
}

// ResourceNotFound 资源不存在
type ResourceNotFound struct {
	Code int    `json:"code" example:"10006"`
	Msg  string `json:"msg" example:"资源不存在"`
}

// TooManyRequest 请求过于频繁
type TooManyRequest struct {
	Code int    `json:"code" example:"10007"`
	Msg  string `json:"msg" example:"请求过于频繁"`
}
