// 流行的web漫画服务xkcd也提供了JSON接口。例如，一个 https://xkcd.com/571/info.0.json 请求将返回一个很多人喜爱的571编号的详细描述。
// 下载每个链接（只下载一次）然后创建一个离线索引。
// 编写一个xkcd工具，使用这些离线索引，打印和命令行输入的检索词相匹配的漫画的URL。
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type XDCD struct {
	Img string
}

func main() {
	result := make(map[int]string)

	for i := 1; i < 1000; i++ {
		uri := "https://xkcd.com/" + strconv.Itoa(i) + "/info.0.json"
		res, err := http.Get(uri)
		if err != nil {
			continue
		}
		var item XDCD
		err = json.NewDecoder(res.Body).Decode(&item)
		if err != nil {
			_ = res.Body.Close()
			continue
		}
		result[i] = item.Img
	}
	_, fullFilename, _, _ := runtime.Caller(0)
	path := filepath.Dir(fullFilename)
	file, err := os.OpenFile(path + "/text.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(file, "%s\t%s\t%s\n", "编号", "访问地址", "漫画地址")
	for index, imgUrl := range result {
		_, _ = fmt.Fprintf(file, "%d\t%s\t%s\n", index, "https://xkcd.com/"+strconv.Itoa(index)+"/info.0.json", imgUrl)
	}
	_ = file.Close()
}
