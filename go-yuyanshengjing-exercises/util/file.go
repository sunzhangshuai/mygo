package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// NewFile 获取新文件
func NewFile(depth int, fName string) (*os.File, error) {
	fileName := GetFileName(depth+1, fName)
	_, err := os.Stat(fileName)
	if os.IsExist(err) {
		if err := os.Remove(fileName); err != nil {
			return nil, err
		}
	}
	file, _ := os.Create(fileName)
	return file, nil
}

// GetFile 获取文件
func GetFile(depth int, fName string) (*os.File, error) {
	fileName := GetFileName(depth+1, fName)
	return os.Open(fileName)
}

// GetFileName 获取文件名称
func GetFileName(depth int, fName string) string {
	_, fullName, _, _ := runtime.Caller(depth)
	return filepath.Join(filepath.Dir(fullName), fName)
}
