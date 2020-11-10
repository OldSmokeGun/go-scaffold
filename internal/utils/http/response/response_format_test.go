package response

import "testing"

func TestFormat(t *testing.T) {
	excepts := map[string]Schema{
		"success_format": {SuccessCode, SuccessCodeMessage, map[string]interface{}{"name": "李四", "age": 18}},
		"failed_format":  {FailedCode, FailedCodeMessage, map[string]interface{}{}},
		"normal_format":   {IllegalRequestCode, IllegalRequestCodeMessage, map[string]interface{}{}},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			schema := Format(v.Code, v.Msg, v.Data)
			if schema.Code != v.Code {
				t.Errorf("http 响应 schema 不匹配，期待 Code：%s，实际 Code：%s", v.Code, schema.Code)
			}
			if schema.Msg != v.Msg {
				t.Errorf("http 响应 schema 不匹配，期待 Msg：%s，实际 Msg：%s", v.Msg, schema.Msg)
			}
			if schema.Data["name"] != v.Data["name"] || schema.Data["age"] != v.Data["age"] {
				t.Errorf("http 响应 schema 不匹配，期待 Data：%v，实际 Data：%v", v.Data, schema.Data)
			}
		})
	}
}

func TestSuccessFormat(t *testing.T) {
	exceptData := map[string]interface{}{"name": "李四", "age": 18}

	schema := SuccessFormat(exceptData)
	if schema.Code != SuccessCode {
		t.Errorf("http 响应 schema 不匹配，期待 Code：%s，实际 Code：%s", SuccessCode, schema.Code)
	}
	if schema.Msg != SuccessCodeMessage {
		t.Errorf("http 响应 schema 不匹配，期待 Msg：%s，实际 Msg：%s", SuccessCodeMessage, schema.Msg)
	}
	if schema.Data["name"] != exceptData["name"] || schema.Data["age"] != exceptData["age"] {
		t.Errorf("http 响应 schema 不匹配，期待 Data：%v，实际 Data：%v", exceptData, schema.Data)
	}
}

func TestFailedFormat(t *testing.T) {
	exceptMsg := "测试错误信息"

	schema := FailedFormat(exceptMsg)
	if schema.Code != FailedCode {
		t.Errorf("http 响应 schema 不匹配，期待 Code：%s，实际 Code：%s", FailedCode, schema.Code)
	}
	if schema.Msg != exceptMsg {
		t.Errorf("http 响应 schema 不匹配，期待 Msg：%s，实际 Msg：%s", exceptMsg, schema.Msg)
	}
	if len(schema.Data) > 0 {
		t.Errorf("http 响应 schema 不匹配，期待 Data：%v，实际 Data：%v", map[string]interface{}{}, schema.Data)
	}
}
