package path

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// RootPath 获取此项目的绝对路径
// 如果是以 go build 生成的二进制文件运行，则返回 bin 目录的上级目录的绝对路径
// 如果是以 go run 运行，则返回在此项目的绝对路径
func RootPath() string {
	var binDir string

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	binDir = filepath.Dir(filepath.Dir(exePath))

	tmpDir := os.TempDir()
	if strings.Contains(exePath, tmpDir) {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			binDir = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
		}
	}

	return binDir
}
