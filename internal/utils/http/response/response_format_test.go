package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	excepts := map[string]Schema{
		"success_format": {SuccessCode, SuccessCodeMessage, map[string]interface{}{"name": "李四", "age": 18}},
		"failed_format":  {FailedCode, FailedCodeMessage, map[string]interface{}{}},
		"other_format":   {IllegalRequestCode, IllegalRequestCodeMessage, map[string]interface{}{}},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			schema := Format(v.Code, v.Msg, v.Data)
			assert.Equal(t, v, schema)
		})
	}
}

func TestSuccessFormat(t *testing.T) {
	excepts := map[string]Schema{
		"data_is_nil":     {SuccessCode, SuccessCodeMessage, map[string]interface{}{}},
		"data_is_not_nil": {SuccessCode, SuccessCodeMessage, map[string]interface{}{"name": "李四", "age": 18}},
	}

	var schema Schema
	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			if k == "data_is_nil" {
				schema = SuccessFormat(nil)
			} else {
				schema = SuccessFormat(v.Data)
			}
			assert.Equal(t, v, schema)
		})
	}
}

func TestFailedFormat(t *testing.T) {
	excepts := map[string]Schema{
		"error_message_is_empty":     {FailedCode, "", map[string]interface{}{}},
		"error_message_is_not_empty": {FailedCode, "测试错误信息", map[string]interface{}{}},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			schema := FailedFormat(v.Msg)
			assert.Equal(t, v, schema)
		})
	}
}
