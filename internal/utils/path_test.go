package utils

import (
	"os"
	"testing"
)

func TestIsDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	testFile, err := os.Executable()
	if err != nil {
		t.Error(err)
	}

	excepts := map[string]map[string]interface{}{
		"is_dir(rel)":     {"path": "../app", "except": true},
		"is_not_dir(rel)": {"path": "../bootstrap.go", "except": false},
		"is_dir(abs)":     {"path": wd, "except": true},
		"is_not_dir(abs)": {"path": testFile, "except": false},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			ok, err := IsDir(v["path"].(string))
			if err != nil {
				t.Errorf("判断 %s 是否为目录出错，错误信息：%s", v["path"].(string), err)
			}
			if ok != v["except"].(bool) {
				t.Errorf("判断 %s 是否为目录失败，期待结果：%t，实际结果：%t", v["path"].(string), v["except"], ok)
			}
		})
	}
}

func TestPathExist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	excepts := map[string]map[string]interface{}{
		"rel_is_exist":     {"path": "../app", "except": true},
		"rel_is_not_exist": {"path": "../not_exist.go", "except": false},
		"abs_is_exist":     {"path": wd, "except": true},
		"abs_is_not_exist": {"path": "/a/b/c/d/e/f", "except": false},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			ok := PathExist(v["path"].(string))
			if ok != v["except"].(bool) {
				t.Errorf("判断路径 %s 是否存在失败，期待结果：%t，实际结果：%t", v["path"].(string), v["except"], ok)
			}
		})
	}
}
