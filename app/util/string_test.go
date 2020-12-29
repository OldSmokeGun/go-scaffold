package util

import "testing"

func TestRandomString(t *testing.T) {
	excepts := map[string]map[string]int{
		"length_0":  {"except": 0},
		"length_32": {"except": 32},
		"length_64": {"except": 64},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			s := RandomString(v["except"])
			if len(s) != v["except"] {
				t.Errorf("生成随机字符串出错，期待字符串长度：%d，实际字符串长度：%d", v["except"], len(s))
			}
		})
	}
}
