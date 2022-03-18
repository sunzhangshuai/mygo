package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type compress struct {
	encoder func(fileOrDir string) error
	decoder func(file string) error
}

// BasePath 基础路径
var BasePath string

var compresses = make(map[string]*compress)

// Register 注册
func Register(action string, encoder, decoder func(fileOrDir string) error) {
	compresses[action] = &compress{
		encoder: encoder,
		decoder: decoder,
	}
}

// Encode 加密
func Encode(action string, fileOrDir string) error {
	if comp, ok := compresses[action]; ok {
		return comp.encoder(fileOrDir)
	}
	return fmt.Errorf("暂不支持%s加密方式", action)
}

// Decode 解压
func Decode(action string, fileOrDir string) error {
	if comp, ok := compresses[action]; ok {
		return comp.decoder(fileOrDir)
	}
	return fmt.Errorf("暂不支持%s解压方式", action)
}

// ExistDir 判断目录是否存在
func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

// init 初始化
func init() {
	// 初始化变量
	_, fullName, _, _ := runtime.Caller(0)
	BasePath = filepath.Join(filepath.Dir(fullName), "file")
}
