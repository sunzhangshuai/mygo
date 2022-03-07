// 编写一个非递归版本的comma函数，使用bytes.Buffer代替字符串链接操作。
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(comma("12321321321"))
	fmt.Println(comma("-12321321321"))
	fmt.Println(comma("+12321321321"))
	fmt.Println(comma("-12321321321.1234"))
	fmt.Println(comma(".123213213211234"))
	fmt.Println(comma("123213213211234."))
}

func comma(s string) string {
	var buf bytes.Buffer
	if s[:1] == "+" || s[:1] == "-" {
		if s[:1] == "-" {
			buf.WriteString(s[:1])
		}
		s = s[1:]
	}
	arr := strings.Split(s, ".")

	// 获取前段
	s = arr[0]
	strLength := len(s)
	if strLength == 0 {
		buf.WriteString("0")
	} else {
		for i := strLength; i > 0; i-- {
			buf.WriteString(s[strLength-i:strLength-i+1])
			if i != 1 && (i - 1) % 3 == 0 {
				buf.WriteString(",")
			}
		}
	}

	// 添加后半段
	if len(arr) == 2 && len(arr[1]) > 0 {
		buf.WriteString(".")
		buf.WriteString(arr[1])
	}
	return buf.String()
}
