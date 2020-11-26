package response

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
	if data == nil {
		data = map[string]interface{}{}
	}

	format := Schema{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	return format
}
