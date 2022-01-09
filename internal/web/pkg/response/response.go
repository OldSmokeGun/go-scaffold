package response

import "reflect"

// Response 响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewResponse 返回 *Response
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// WithCode 设置响应码
func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

// WithMsg 设置响应消息
func (r *Response) WithMsg(msg string) *Response {
	r.Msg = msg
	return r
}

// WithData 设置响应数据
func (r *Response) WithData(data interface{}) *Response {
	if data != nil {
		val := reflect.ValueOf(data)
		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		case reflect.Slice:
			if val.IsNil() {
				data = make([]interface{}, 0)
			}
		case reflect.Map:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		case reflect.Chan:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		case reflect.Func:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		case reflect.Interface:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		}
	} else {
		data = map[string]interface{}{}
	}

	r.Data = data

	return r
}
