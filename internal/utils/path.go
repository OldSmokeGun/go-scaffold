package utils

import (
	"fmt"
	"os"
)

func BasePath() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径出错，错误信息：%w", err)
	}
	return p, nil
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsExist(err) {
		return true, err
	}

	return false, err
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
