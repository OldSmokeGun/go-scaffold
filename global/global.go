package global

import (
	"os"
)

var (
	err     error
	binPath string // 二进制文件路径
)

func init() {
	binPath, err = os.Executable()
	if err != nil {
		panic(err)
	}
}

// BinPath 获取二进制文件路径
func GetBinPath() string {
	return binPath
}
