package response

type Schema struct {
	Code string                 `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func SuccessFormat(data map[string]interface{}) Schema {
	return Format(SuccessCode, SuccessCodeMessage, data)
}

func FailedFormat(msg string) Schema {
	return Format(FailedCode, msg, map[string]interface{}{})
}

func Format(code string, msg string, data map[string]interface{}) Schema {
	format := Schema{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	return format
}
