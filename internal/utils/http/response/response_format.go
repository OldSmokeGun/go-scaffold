package response

import "reflect"

type Schema struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessFormat(data interface{}) Schema {
	return Format(SuccessCode, SuccessCodeMessage, data)
}

func FailedFormat(msg string) Schema {
	return Format(FailedCode, msg, map[string]interface{}{})
}

func Format(code string, msg string, data interface{}) Schema {
	if data != nil {
		val := reflect.ValueOf(data)
		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				data = map[string]interface{}{}
			}
		case reflect.Slice:
			if val.IsNil() {
				data = map[string]interface{}{}
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

	format := Schema{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	return format
}
