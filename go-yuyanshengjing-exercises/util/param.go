package util

import (
	"os"
)

// GetParam 获取参数
func GetParam() []string {
	return os.Args[3:]
}
