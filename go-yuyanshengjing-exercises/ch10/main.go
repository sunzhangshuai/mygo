package ch10

import (
	_ "archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"exercises/ch10/compress"
	_ "exercises/ch10/compress/tar"
	_ "exercises/ch10/compress/zip"
	"exercises/ch9"
	"exercises/util"
)

// Exercises 练习
type Exercises struct {
}

// trans 图片转换
func trans(in io.Reader, out io.Writer, typ string) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	switch typ {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{})
	}
	return fmt.Errorf("typ error")
}

// Task1 扩展jpeg程序，以支持任意图像格式之间的相互转换，使用image.Decode检测支持的格式类型，然后通过flag命令行标志参数选择输出的格式。
// go run main.go -ch 10 -task 1 jpeg <./ch10/file/image.gif >./ch10/file/image.jpeg
func (e *Exercises) Task1() {
	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("参数不足")
		return
	}
	if err := trans(os.Stdin, os.Stdout, params[0]); err != nil {
		fmt.Printf("jpeg: %v\n", err)
		os.Exit(1)
	}
}

// Task2 设计一个通用的压缩文件读取框架，用来读取ZIP（archive/zip）和POSIX tar（archive/tar）格式压缩的文档。
// 使用类似上面的注册技术来扩展支持不同的压缩格式，然后根据需要通过匿名导入选择导入要支持的压缩格式的驱动包。
// go run main.go -ch 10 -task 2 zip encode test
// go run main.go -ch 10 -task 2 zip decode test.zip
// go run main.go -ch 10 -task 2 tar encode test
// go run main.go -ch 10 -task 2 tar decode test.tar.gz
func (e *Exercises) Task2() {
	params := util.GetParam()
	if len(params) == 3 {
		fmt.Println("参数不足")
		return
	}
	action := params[1] // 压缩算法
	typ := params[2]    // 解压缩or压缩
	filename := params[3]

	_, fullName, _, _ := runtime.Caller(0)
	filePath := path.Dir(fullName)
	filename = path.Join(filePath, "compress", "file", filename)

	switch typ {
	case "decode":
		fmt.Println(compress.Decode(action, filename))
	case "encode":
		fmt.Println(compress.Encode(action, filename))
	}
}

// Task3 从 http://gopl.io/ch1/helloworld?go-get=1 获取内容，查看本书的代码的真实托管的网址（go get请求HTML页面时包含了go-get参数，以区别普通的浏览器请求）。
// go run main.go -ch=10 -task=3 zhangshuai
func (e *Exercises) Task3() {
	var err error
	var file *os.File

	params := util.GetParam()
	if len(params) == 0 {
		fmt.Println("请输入文件名")
	}

	if file, err = util.NewFile(1, filepath.Join("file", params[0]+".html")); err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	res, _ := http.Get("http://gopl.io/ch1/helloworld?go-get=1")
	defer res.Body.Close()
	out := io.MultiWriter(file, os.Stdout)
	io.Copy(out, res.Body)
}

// Task4 创建一个工具，根据命令行指定的参数，报告工作区所有依赖包指定的其它包集合。
// 提示：你需要运行go list命令两次，一次用于初始化包，一次用于所有包。
// 你可能需要用encoding/json（§4.5）包来分析输出的JSON格式的信息。
// go run main.go -ch 10 -task 4
func (e *Exercises) Task4() {
	cmd := exec.Command("go", "list", "./...")
	res, _ := cmd.Output()
	buf := bytes.NewBuffer(res)
	scaner := bufio.NewScanner(buf)

	type pkgInfo struct {
		Imports []string `json:"Imports"`
	}

	var count int
	memo := ch9.NewMemo(func(key string, done chan struct{}) (interface{}, error) {
		count++
		return true, nil
	})

	group := sync.WaitGroup{}

	for scaner.Scan() {
		pkgName := scaner.Text()
		group.Add(1)
		go func() {
			defer group.Done()
			cmd := exec.Command("go", "list", "-json", pkgName)
			res, _ := cmd.Output()
			resMap := &pkgInfo{}
			json.Unmarshal(res, resMap)
			for _, val := range resMap.Imports {
				group.Add(1)
				go func(val string) {
					defer group.Done()
					memo.Get(val, nil)
				}(val)
			}
		}()
	}

	group.Wait()
	fmt.Println(memo.GetAllKey(), len(memo.GetAllKey()), count)
}
