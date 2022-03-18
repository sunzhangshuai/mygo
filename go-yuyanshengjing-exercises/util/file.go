package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// NewFile 获取新文件
func NewFile(depth int, fName string) (*os.File, error) {
	_, fullName, _, _ := runtime.Caller(depth)
	fmt.Println(fullName)
	fileName := filepath.Join(filepath.Dir(fullName), fName)

	_, err := os.Stat(fileName)

	if os.IsExist(err) {
		if err := os.Remove(fileName); err != nil {
			return nil, err
		}
	}

	file, _ := os.Create(fileName)
	return file, nil
}
